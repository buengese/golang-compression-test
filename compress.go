package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	flstd "compress/flate"
	gzstd "compress/gzip"

	bgzf "github.com/biogo/hts/bgzf"
	flkp "github.com/klauspost/compress/flate"
	gzkp "github.com/klauspost/compress/gzip"
	zstd "github.com/klauspost/compress/zstd"
	pgz "github.com/klauspost/pgzip"
	lz4 "github.com/pierrec/lz4"
	gxz "github.com/ulikunitz/xz/lzma"

	clz4 "github.com/pwaller/go-clz4"
	clzma "github.com/remyoudompheng/go-liblzma"
	cgzip "github.com/youtube/vitess/go/cgzip"

	"github.com/golang/snappy"
	"github.com/klauspost/readahead"
)

type NoOp struct{}

func (n NoOp) Read(v []byte) (int, error) {
	return len(v), nil
}

func (n NoOp) Write(v []byte) (int, error) {
	return len(v), nil
}

type SeqGen struct {
	i int
}

func (s *SeqGen) Read(v []byte) (int, error) {
	b := byte(s.i)
	for i := range v {
		v[i], b = b, b+1
	}
	return len(v), nil
}

type Rand struct {
	state uint64
	inc   uint64
}

const pcgmult64 = 6364136223846793005

func NewRand(seed uint64) *Rand {
	state := uint64(0)
	inc := uint64(seed<<1) | 1
	state = state*pcgmult64 + (inc | 1)
	state += uint64(seed)
	state = state*pcgmult64 + (inc | 1)
	return &Rand{
		state: state,
		inc:   inc,
	}
}

func (r *Rand) Read(v []byte) (int, error) {
	for w := v; len(w) > 0; w = w[4:] {
		old := r.state
		r.state = r.state*pcgmult64 + (r.inc | 1)
		xorshifted := uint32(((old >> 18) ^ old) >> 27)
		rot := uint32(old >> 59)
		rnd := (xorshifted >> rot) | (xorshifted << ((-rot) & 31))
		// ok because len(v) % 4 == 0
		binary.LittleEndian.PutUint32(w, rnd)
	}
	return len(v), nil
}

type wcounter struct {
	n   int
	out io.Writer
}

func (w *wcounter) Write(p []byte) (n int, err error) {
	n, err = w.out.Write(p)
	w.n += n
	return n, err

}

func main() {
	rmode := "raw"
	wmode := "gzkp"
	wlevel := -1
	in := "-"
	out := "-"
	cpu := 0
	stats := false
	header := true

	flag.StringVar(&rmode, "r", rmode, "read mode (raw|flatekp|flatestd|gzkp|pgzip|cgzip|gzstd|zero|seq|rand)")
	flag.StringVar(&wmode, "w", wmode, "write mode (raw|flatekp|flatestd|gzkp|pgzip|gzstd|cgzip|none)")
	flag.StringVar(&in, "in", rmode, "input file name, default is '-', stdin")
	flag.StringVar(&out, "out", rmode, "input file name, default is '-', stdout")
	flag.IntVar(&wlevel, "l", wlevel, "compression level (-2|-1|0..9)")
	flag.IntVar(&cpu, "cpu", cpu, "GOMAXPROCS number (0|1...)")
	flag.BoolVar(&stats, "stats", false, "show stats")
	flag.BoolVar(&header, "header", true, "show stats header")
	flag.Parse()
	if flag.NArg() > 0 {
		flag.PrintDefaults()
	}
	if cpu <= 0 {
		cpu = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(cpu)

	if wlevel < -2 || 9 < wlevel {
		panic("compression level -l=x must be (-2,0..9)")
	}

	var err error

	var r io.Reader
	if in == "-" {
		r = os.Stdin
	} else {
		r, err = os.Open(in)
		if err != nil {
			panic(err)
		}
		r, _ = readahead.NewReaderSize(r, 10, 10<<20)
	}
	var source bool
	switch rmode {
	case "zero":
		// NoOp writes what the original buffer contained unchanged.
		// As that buffer is initialized with 0 and not changed,
		// NoOp is usable as a very fast zero-reader.
		r = NoOp{}
		source = true
	case "seq":
		r = &SeqGen{}
		source = true
	case "rand":
		r = NewRand(0xdeadbeef)
		source = true
	case "raw":
	case "gzkp":
		var gzr *gzkp.Reader
		if gzr, err = gzkp.NewReader(r); err == nil {
			defer gzr.Close()
			r = gzr
		}
	case "bgzf":
		var gzr *bgzf.Reader
		if gzr, err = bgzf.NewReader(r, cpu); err == nil {
			defer gzr.Close()
			r = gzr
		}
	case "pgzip":
		var gzr *pgz.Reader
		if gzr, err = pgz.NewReader(r); err == nil {
			defer gzr.Close()
			r = gzr
		}
	case "cgzip":
		var gzr io.ReadCloser
		if gzr, err = cgzip.NewReader(r); err == nil {
			defer gzr.Close()
			r = gzr
		}
	case "gzstd":
		var gzr *gzstd.Reader
		if gzr, err = gzstd.NewReader(r); err == nil {
			defer gzr.Close()
			r = gzr
		}
	case "snappy":
		sr := snappy.NewReader(r)
		r = sr
	case "flatekp":
		fr := flkp.NewReader(r)
		defer fr.Close()
		r = fr
	case "flatestd":
		fr := flstd.NewReader(r)
		defer fr.Close()
		r = fr
	default:
		panic("read mode -r=x must be (raw|flatekp|flatestd|gzkp|gzstd|zero|seq|rand)")
	}
	if err != nil {
		panic(err)
	}

	var w io.Writer
	if out == "-" {
		w = os.Stdout
	} else if out == "*" {
		w = ioutil.Discard
		out = "discard"
	} else {
		f, err := os.Create(out)
		if err != nil {
			panic(err)
		}
		w = bufio.NewWriter(f)
	}
	outSize := &wcounter{out: w}
	w = outSize

	var sink bool
	switch wmode {
	case "none":
		w = NoOp{}
		sink = true
	case "raw":
	case "gzkp":
		cpu = 1
		var gzw *gzkp.Writer
		if gzw, err = gzkp.NewWriterLevel(w, wlevel); err == nil {
			defer gzw.Close()
			w = gzw
		}
	case "pgzip":
		var gzw *pgz.Writer
		if gzw, err = pgz.NewWriterLevel(w, wlevel); err == nil {
			defer gzw.Close()
			w = gzw
		}
	case "bgzf":
		var gzw *bgzf.Writer
		if gzw, err = bgzf.NewWriterLevel(w, wlevel, cpu); err == nil {
			defer gzw.Close()
			w = gzw
		}
	case "cgzip":
		cpu = 1
		var gzw *cgzip.Writer
		if gzw, err = cgzip.NewWriterLevel(w, wlevel); err == nil {
			defer gzw.Close()
			w = gzw
		}
	case "gzstd":
		cpu = 1
		var gzw *gzstd.Writer
		if gzw, err = gzstd.NewWriterLevel(w, wlevel); err == nil {
			defer gzw.Close()
			w = gzw
		}
	case "flatekp":
		cpu = 1
		var fw *flkp.Writer
		if fw, err = flkp.NewWriter(w, wlevel); err == nil {
			defer fw.Close()
			w = fw
		}
	case "flatestd":
		cpu = 1
		var fw *flstd.Writer
		if fw, err = flstd.NewWriter(w, wlevel); err == nil {
			defer fw.Close()
			w = fw
		}
	case "clzma":
		cpu = 1
		var lzmaw *clzma.Compressor
		if lzmaw, err = clzma.NewWriter(w, clzma.Preset(wlevel)); err == nil {
			w = lzmaw
		}
	case "gxz":
		cpu = 1
		var lzmaw *gxz.Writer
		if lzmaw, err = gxz.NewWriter(w); err == nil {
			defer lzmaw.Close()
			w = lzmaw
		}
	case "gxz2":
		cpu = 1
		var lzmaw *gxz.Writer2
		if lzmaw, err = gxz.NewWriter2(w); err == nil {
			defer lzmaw.Close()
			w = lzmaw
		}
	case "snappy":
		cpu = 1
		w = snappy.NewWriter(w)
	case "lz4":
		cpu = 1
		lzw := lz4.NewWriter(w)
		defer lzw.Close()
		w = lzw
	case "lz4p":
		lzw := lz4.NewWriter(w).WithConcurrency(cpu)
		defer lzw.Close()
		w = lzw
	case "clz4":
		lzw := clz4.NewWriter(w)
		w = lzw
	case "zstdfast":
		cpu = 2
		var zstdw *zstd.Encoder
		if zstdw, err = zstd.NewWriter(w, zstd.WithEncoderLevel(zstd.SpeedFastest), zstd.WithEncoderConcurrency(2)); err == nil {
			defer zstdw.Close()
			w = zstdw
		}
	case "zstdbest":
		cpu = 2
		var zstdw *zstd.Encoder
		if zstdw, err = zstd.NewWriter(w, zstd.WithEncoderLevel(zstd.SpeedDefault), zstd.WithEncoderConcurrency(2)); err == nil {
			defer zstdw.Close()
			w = zstdw
		}
	default:
		panic("write mode -w=x must be (raw|gzkp|pgzip|bgzf|cgzip|gzstd|flatekp|flatestd|clzma|gxz|gxz2|snappy|lz4|lz4p|zsdfast|zstdbest|none)")
	}
	if err != nil {
		panic(err)
	}

	if source && sink {
		return
	}

	inSize := int64(0)
	start := time.Now()
	func() {
		for _, mc := range []interface{}{r, w} {
			if c, ok := mc.(io.Closer); ok {
				defer c.Close()
			}
		}

		nr, err := io.Copy(w, r)
		inSize += nr
		if err != nil && err != io.EOF {
			panic(err)
		}
	}()

	if stats {
		elapsed := time.Since(start)
		if header {
			fmt.Printf("file\tin\tout\tlevel\tcpu\tinsize\toutsize\tmillis\tmb/s\n")
		}
		mbpersec := (float64(inSize) / (1024 * 1024)) / (float64(elapsed) / (float64(time.Second)))
		fmt.Printf("%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%.02f\n", in, rmode, wmode, wlevel, cpu, inSize, outSize.n, elapsed/time.Millisecond, mbpersec)
	}
}
