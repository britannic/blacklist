# regx
--
    import "github.com/britannic/blacklist/internal/regx"

Package regx provides regex objects for processing data in files and web content

## Usage

#### type Leaf

```go
type Leaf int
```

Leaf is a config label

```go
const (
	CMNT Leaf = iota + 1000
	DESC
	DSBL
	FLIP
	FQDN
	HOST
	HTTP
	IPBH
	LEAF
	LBRC
	MISC
	MLTI
	MPTY
	NAME
	NODE
	RBRC
	SUFX
)
```
go:generate stringer -type=Leaf Leaf label regx map keys

#### func (Leaf) String

```go
func (i Leaf) String() string
```

#### type OBJ

```go
type OBJ struct {
	RX regexMap
}
```

OBJ is a struct of regex precompiled objects

#### func  NewRegex

```go
func NewRegex() *OBJ
```
NewRegex returns a map of OBJ populated with a map of precompiled regex objects

#### func (*OBJ) String

```go
func (o *OBJ) String() string
```

#### func (*OBJ) StripPrefixAndSuffix

```go
func (o *OBJ) StripPrefixAndSuffix(l []byte, p string) ([]byte, bool)
```
StripPrefixAndSuffix strips the prefix and suffix

#### func (*OBJ) SubMatch

```go
func (o *OBJ) SubMatch(t Leaf, b []byte) [][]byte
```
SubMatch extracts the configuration value for a matched label
