
# edgeos
    import "github.com/britannic/blacklist/internal/edgeos"

Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.




## Constants
``` go
const (

    // False is a string constant
    False = "false"
    // True is a string constant
    True = "true"
)
```


## func BooltoStr
``` go
func BooltoStr(b bool) string
```
BooltoStr converts a boolean ("true" or "false") to a string equivalent


## func DiffArray
``` go
func DiffArray(a, b []string) (diff sort.StringSlice)
```
DiffArray returns the delta of two arrays


## func GetHTTP
``` go
func GetHTTP(method, URL string) (io.Reader, error)
```
GetHTTP creates http requests to download data


## func StrToBool
``` go
func StrToBool(s string) bool
```
StrToBool converts a string ("true" or "false") to it's boolean equivalent



## type CFGcli
``` go
type CFGcli struct {
    Cfg string
}
```
CFGcli is for configurations loaded via the EdgeOS CFGcli











### func (\*CFGcli) Load
``` go
func (c *CFGcli) Load() io.Reader
```
Load returns an EdgeOS config file string and error



## type CFGstatic
``` go
type CFGstatic struct {
    Cfg string
}
```
CFGstatic is for configurations loaded via the EdgeOS CFGstatic











### func (\*CFGstatic) Load
``` go
func (c *CFGstatic) Load() io.Reader
```
Load returns an EdgeOS CLI loaded configuration



## type CFile
``` go
type CFile struct {
    *Parms
    // contains filtered or unexported fields
}
```
CFile holds an array of file names











### func (\*CFile) ReadDir
``` go
func (c *CFile) ReadDir(dir string) ([]os.FileInfo, error)
```
ReadDir implements OSinformer



### func (\*CFile) Remove
``` go
func (c *CFile) Remove() error
```
Remove deletes a CFile array of file names



### func (\*CFile) String
``` go
func (c *CFile) String() string
```
String implements string method



### func (\*CFile) Strings
``` go
func (c *CFile) Strings() []string
```
Strings returns a sorted array of strings.



## type ConfLoader
``` go
type ConfLoader interface {
    Load() io.Reader
}
```
ConfLoader interface defines load method











## type Config
``` go
type Config struct {
    *Parms
    // contains filtered or unexported fields
}
```
Config is a struct of configuration fields









### func ReadCfg
``` go
func ReadCfg(r ConfLoader) (*Config, error)
```
ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure




### func (\*Config) Excludes
``` go
func (c *Config) Excludes(node string) []string
```
Excludes returns a string array of excludes



### func (\*Config) Get
``` go
func (c *Config) Get(node string) (o *Object)
```
Get returns an *Object for a given node



### func (\*Config) GetAll
``` go
func (c *Config) GetAll() *Objects
```
GetAll returns an array of Objects



### func (\*Config) Nodes
``` go
func (c *Config) Nodes() (nodes []string)
```
Nodes returns an array of configured nodes



### func (\*Config) STypes
``` go
func (c *Config) STypes() []string
```
STypes returns an array of configured nodes



### func (\*Config) String
``` go
func (c *Config) String() (result string)
```
String returns pretty print for the Blacklist struct



## type Content
``` go
type Content struct {
    *Object
    Contenter

    *Parms
    // contains filtered or unexported fields
}
```
Content is a struct of blacklist content











### func (\*Content) Process
``` go
func (c *Content) Process() *blist
```
Process extracts hosts/domains from downloaded raw content



### func (Content) Source
``` go
func (d Content) Source(ltype string) *Objects
```
Source returns a map of sources



## type Contenter
``` go
type Contenter interface {
    Process() io.Reader
}
```
Contenter is a Content interface











## type Contents
``` go
type Contents []*Content
```
Contents is an array of *content











### func (\*Contents) ProcessContent
``` go
func (c *Contents) ProcessContent()
```
ProcessContent iterates through the Contents array and processes each



### func (\*Contents) String
``` go
func (c *Contents) String() (result string)
```


## type List
``` go
type List map[string]int
```
List is a map of int









### func UpdateList
``` go
func UpdateList(data []string) (l List)
```
UpdateList converts []string to map of List




### func (List) String
``` go
func (l List) String() string
```
String implements fmt.Print interface



## type OSinformer
``` go
type OSinformer interface {
    ReadDir(string) ([]os.FileInfo, error)
    Remove() error
}
```
OSinformer implements os.FileInfo methods











## type Object
``` go
type Object struct {
    *Parms
    // contains filtered or unexported fields
}
```
Object struct for normalizing EdgeOS data.











### func (\*Object) Excludes
``` go
func (o *Object) Excludes() List
```
Excludes returns a List map of blacklist exclusions



### func (\*Object) Includes
``` go
func (o *Object) Includes() io.Reader
```
Includes returns an io.Reader of blacklist Includes



### func (Object) Source
``` go
func (d Object) Source(ltype string) *Objects
```
Source returns a map of sources



### func (\*Object) String
``` go
func (o *Object) String() (r string)
```
Stringer for Object



## type Objects
``` go
type Objects struct {
    S []*Object
    *Parms
}
```
Objects is a struct of []*Object











### func (\*Objects) Files
``` go
func (o *Objects) Files() *CFile
```
Files returns a list of dnsmasq conf files from all srcs



### func (\*Objects) GetContent
``` go
func (objs *Objects) GetContent() *Contents
```
GetContent returns a Content struct



### func (\*Objects) Len
``` go
func (o *Objects) Len() int
```
Implement Sort Interface for Objects



### func (\*Objects) Less
``` go
func (o *Objects) Less(i, j int) bool
```


### func (\*Objects) String
``` go
func (o *Objects) String() string
```
Stringer for Objects



### func (\*Objects) Swap
``` go
func (o *Objects) Swap(i, j int)
```


## type Option
``` go
type Option func(p *Parms) Option
```
Option sets is a recursive function









### func Arch
``` go
func Arch(arch string) Option
```
Arch sets target CPU architecture


### func Cores
``` go
func Cores(i int) Option
```
Cores sets max CPU cores


### func Debug
``` go
func Debug(b bool) Option
```
Debug toggles debug level on or off


### func Dir
``` go
func Dir(d string) Option
```
Dir sets directory location


### func Excludes
``` go
func Excludes(l List) Option
```
Excludes sets nodes exclusions


### func Ext
``` go
func Ext(e string) Option
```
Ext sets the blacklist file n extension


### func File
``` go
func File(f string) Option
```
File sets the EdgeOS configuration file n


### func FileNameFmt
``` go
func FileNameFmt(f string) Option
```
FileNameFmt sets the EdgeOS configuration file name format


### func Method
``` go
func Method(method string) Option
```
Method sets the HTTP method


### func Nodes
``` go
func Nodes(nodes []string) Option
```
Nodes sets the node ns array


### func Poll
``` go
func Poll(t int) Option
```
Poll sets the polling interval in seconds


### func Prefix
``` go
func Prefix(l string) Option
```
Prefix sets the dnsmasq configuration address line prefix


### func STypes
``` go
func STypes(s []string) Option
```
STypes sets an array of legal types used by Source


### func Test
``` go
func Test(b bool) Option
```
Test toggles testing mode on or off


### func Verbosity
``` go
func Verbosity(i int) Option
```
Verbosity sets the verbosity level to v




## type Parms
``` go
type Parms struct {
    Arch      string
    Cores     int
    Debug     bool
    Dex       List
    Dir       string
    Exc       List
    Ext       string
    File      string
    FnFmt     string
    Method    string
    Nodes     []string
    Pfx       string
    Poll      int
    Stypes    []string
    Test      bool
    Verbosity int
}
```
Parms is struct of parameters









### func NewParms
``` go
func NewParms() *Parms
```
NewParms sets a new *Parms instance




### func (\*Parms) SetOpt
``` go
func (p *Parms) SetOpt(opts ...Option) (previous Option)
```
SetOpt sets the specified options passed as Parms and returns an option to restore the last arg's previous value



### func (\*Parms) String
``` go
func (p *Parms) String() string
```
String method to implement fmt.Print interface









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)