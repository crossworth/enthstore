package enthstore

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// testdata from:
// https://github.com/postgres/postgres/blob/d33a81203e95d31e62157c4ae0e00e2198841208/contrib/hstore/sql/hstore.sql
var testdata = `

a=>b
 a=>b
a =>b
a=>b 
a=> b
"a"=>"b"
 "a"=>"b"
"a" =>"b"
"a"=>"b" 
"a"=> "b"
aa=>bb
 aa=>bb
aa =>bb
aa=>bb 
aa=> bb
"aa"=>"bb"
 "aa"=>"bb"
"aa" =>"bb"
"aa"=>"bb" 
"aa"=> "bb"
aa=>bb, cc=>dd
aa=>bb , cc=>dd
aa=>bb ,cc=>dd
aa=>bb, "cc"=>dd
aa=>bb , "cc"=>dd
aa=>bb ,"cc"=>dd
aa=>"bb", cc=>dd
aa=>"bb" , cc=>dd
aa=>"bb" ,cc=>dd
aa=>null
aa=>NuLl
aa=>"NuLl"
\\=a=>q=w
"=a"=>q\\=w
"\\"a"=>q>w
\\"a=>q"w

	
`

func Test_canParse(t *testing.T) {
	lines := strings.Split(testdata, "\n")
	for _, l := range lines {
		hs := Hstore{}
		err := hs.Scan(l)
		require.NoError(t, err)
	}
}

func TestHstore_Scan(t *testing.T) {
	tests := []struct {
		args       interface{}
		wantHstore Hstore
		wantErr    error
	}{
		{args: "", wantHstore: nil, wantErr: nil},
		{args: "a=>b", wantHstore: FromMap(map[string]string{"a": "b"}), wantErr: nil},
		{args: " a=>b", wantHstore: FromMap(map[string]string{"a": "b"}), wantErr: nil},
		{args: `a=>b`, wantHstore: FromMap(map[string]string{"a": "b"}), wantErr: nil},
		{args: ` a=>b`, wantHstore: FromMap(map[string]string{"a": "b"}), wantErr: nil},
		{args: `a =>b`, wantHstore: FromMap(map[string]string{"a": "b"}), wantErr: nil},
		{args: `a=>b`, wantHstore: FromMap(map[string]string{"a": "b"}), wantErr: nil},
		{args: `a=> b`, wantHstore: FromMap(map[string]string{"a": "b"}), wantErr: nil},
		{args: `"a"=>"b"`, wantHstore: FromMap(map[string]string{"a": "b"}), wantErr: nil},
		{args: ` "a"=>"b"`, wantHstore: FromMap(map[string]string{"a": "b"}), wantErr: nil},
		{args: `"a" =>"b"`, wantHstore: FromMap(map[string]string{"a": "b"}), wantErr: nil},
		{args: `"a"=>"b"`, wantHstore: FromMap(map[string]string{"a": "b"}), wantErr: nil},
		{args: `"a"=> "b"`, wantHstore: FromMap(map[string]string{"a": "b"}), wantErr: nil},
		{args: `aa=>bb`, wantHstore: FromMap(map[string]string{"aa": "bb"}), wantErr: nil},
		{args: ` aa=>bb`, wantHstore: FromMap(map[string]string{"aa": "bb"}), wantErr: nil},
		{args: `aa =>bb`, wantHstore: FromMap(map[string]string{"aa": "bb"}), wantErr: nil},
		{args: `aa=>bb`, wantHstore: FromMap(map[string]string{"aa": "bb"}), wantErr: nil},
		{args: `aa=> bb`, wantHstore: FromMap(map[string]string{"aa": "bb"}), wantErr: nil},
		{args: `"aa"=>"bb"`, wantHstore: FromMap(map[string]string{"aa": "bb"}), wantErr: nil},
		{args: ` "aa"=>"bb"`, wantHstore: FromMap(map[string]string{"aa": "bb"}), wantErr: nil},
		{args: `"aa" =>"bb"`, wantHstore: FromMap(map[string]string{"aa": "bb"}), wantErr: nil},
		{args: `"aa"=>"bb"`, wantHstore: FromMap(map[string]string{"aa": "bb"}), wantErr: nil},
		{args: `"aa"=> "bb"`, wantHstore: FromMap(map[string]string{"aa": "bb"}), wantErr: nil},
		{args: `aa=>bb, cc=>dd`, wantHstore: FromMap(map[string]string{"aa": "bb", "cc": "dd"}), wantErr: nil},
		{args: `aa=>bb , cc=>dd`, wantHstore: FromMap(map[string]string{"aa": "bb", "cc": "dd"}), wantErr: nil},
		{args: `aa=>bb ,cc=>dd`, wantHstore: FromMap(map[string]string{"aa": "bb", "cc": "dd"}), wantErr: nil},
		{args: `aa=>bb, "cc"=>dd`, wantHstore: FromMap(map[string]string{"aa": "bb", "cc": "dd"}), wantErr: nil},
		{args: `aa=>bb , "cc"=>dd`, wantHstore: FromMap(map[string]string{"aa": "bb", "cc": "dd"}), wantErr: nil},
		{args: `aa=>bb ,"cc"=>dd`, wantHstore: FromMap(map[string]string{"aa": "bb", "cc": "dd"}), wantErr: nil},
		{args: `aa=>"bb", cc=>dd`, wantHstore: FromMap(map[string]string{"aa": "bb", "cc": "dd"}), wantErr: nil},
		{args: `aa=>"bb" , cc=>dd`, wantHstore: FromMap(map[string]string{"aa": "bb", "cc": "dd"}), wantErr: nil},
		{args: `aa=>"bb" ,cc=>dd`, wantHstore: FromMap(map[string]string{"aa": "bb", "cc": "dd"}), wantErr: nil},
		{args: `aa=>null`, wantHstore: Hstore{"aa": nil}, wantErr: nil},
		{args: `aa=>NuLl`, wantHstore: Hstore{"aa": nil}, wantErr: nil},
		{args: `aa=>"NuLl"`, wantHstore: FromMap(map[string]string{"aa": "NuLl"}), wantErr: nil},
		{args: `	`, wantHstore: nil, wantErr: nil},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var hs Hstore
			err := hs.Scan(tt.args)
			require.Equal(t, tt.wantErr, err)
			require.True(t, hs.Equals(tt.wantHstore))
		})
	}
}

func TestHstore_String(t *testing.T) {
	h1 := FromMap(map[string]string{
		"key1": "value1",
		"key2": "value2",
	})

	res1 := h1.String()
	require.True(t, res1 == `"key1"=>"value1","key2"=>"value2"` ||
		res1 == `"key2"=>"value2","key1"=>"value1"`)

	h2 := Hstore{}
	h2.SetString("key1", "value1")
	h2.Set("key2", nil)

	res2 := h2.String()
	require.True(t, res2 == `"key1"=>"value1","key2"=>NULL` ||
		res2 == `"key2"=>NULL,"key1"=>"value1"`)
}

func TestHstore_UnmarshalGQL(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		in := map[string]interface{}{"a": "b"}

		h := &Hstore{}
		err := h.UnmarshalGQL(in)
		require.NoError(t, err)
		require.True(t, h.Has("a"))
		require.Equal(t, "b", h.GetString("a"))
	})

	t.Run("null", func(t *testing.T) {
		in := map[string]interface{}{"a": "b", "c": nil}

		h := &Hstore{}
		err := h.UnmarshalGQL(in)
		require.NoError(t, err)
		require.True(t, h.Has("a"))
		require.Equal(t, "b", h.GetString("a"))
		require.Nil(t, h.Get("c"))
	})
}

func TestHstore_MarshalGQL(t *testing.T) {
	h1 := FromMap(map[string]string{"a": "b", "c": "d"})

	var out bytes.Buffer
	h1.MarshalGQL(&out)
	res1 := out.String()
	require.True(t, res1 == "{\"a\":\"b\",\"c\":\"d\"}\n" || res1 == "{\"c\":\"d\",\"a\":\"b\"}\n")

	out.Reset()

	h2 := Hstore{}
	h2.SetString("a", "b")
	h2.Set("c", nil)
	h2.MarshalGQL(&out)
	res2 := out.String()
	require.True(t, res2 == "{\"a\":\"b\",\"c\":null}\n" || res2 == "{\"c\":null,\"a\":\"b\"}\n")
}
