package gbconv_test

import (
	gbcrc32 "ghostbb.io/gb/crypto/gb_crc32"
	gbbinary "ghostbb.io/gb/encoding/gb_binary"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/frame/g"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
	"time"
)

type MyTime struct {
	time.Time
}

type MyTimeSt struct {
	ServiceDate MyTime
}

func (st *MyTimeSt) UnmarshalValue(v interface{}) error {
	m := gbconv.Map(v)
	t, err := gbtime.StrToTime(gbconv.String(m["ServiceDate"]))
	if err != nil {
		return err
	}
	st.ServiceDate = MyTime{t.Time}
	return nil
}

func Test_Struct_UnmarshalValue1(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		st := &MyTimeSt{}
		err := gbconv.Struct(g.Map{"ServiceDate": "2020-10-10 12:00:01"}, st)
		t.AssertNil(err)
		t.Assert(st.ServiceDate.Time.Format("2006-01-02 15:04:05"), "2020-10-10 12:00:01")
	})
	gbtest.C(t, func(t *gbtest.T) {
		st := &MyTimeSt{}
		err := gbconv.Struct(g.Map{"ServiceDate": nil}, st)
		t.AssertNil(err)
		t.Assert(st.ServiceDate.Time.IsZero(), true)
	})
	gbtest.C(t, func(t *gbtest.T) {
		st := &MyTimeSt{}
		err := gbconv.Struct(g.Map{"ServiceDate": "error"}, st)
		t.AssertNE(err, nil)
	})
}

type Pkg struct {
	Length uint16 // Total length.
	Crc32  uint32 // CRC32.
	Data   []byte
}

// NewPkg creates and returns a package with given data.
func NewPkg(data []byte) *Pkg {
	return &Pkg{
		Length: uint16(len(data) + 6),
		Crc32:  gbcrc32.Encrypt(data),
		Data:   data,
	}
}

// Marshal encodes the protocol struct to bytes.
func (p *Pkg) Marshal() []byte {
	b := make([]byte, 6+len(p.Data))
	copy(b, gbbinary.EncodeUint16(p.Length))
	copy(b[2:], gbbinary.EncodeUint32(p.Crc32))
	copy(b[6:], p.Data)
	return b
}

// UnmarshalValue decodes bytes to protocol struct.
func (p *Pkg) UnmarshalValue(v interface{}) error {
	b := gbconv.Bytes(v)
	if len(b) < 6 {
		return gberror.New("invalid package length")
	}
	p.Length = gbbinary.DecodeToUint16(b[:2])
	if len(b) < int(p.Length) {
		return gberror.New("invalid data length")
	}
	p.Crc32 = gbbinary.DecodeToUint32(b[2:6])
	p.Data = b[6:]
	if gbcrc32.Encrypt(p.Data) != p.Crc32 {
		return gberror.New("crc32 validation failed")
	}
	return nil
}

func Test_Struct_UnmarshalValue2(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var p1, p2 *Pkg
		p1 = NewPkg([]byte("123"))
		err := gbconv.Struct(p1.Marshal(), &p2)
		t.AssertNil(err)
		t.Assert(p1, p2)
	})
}
