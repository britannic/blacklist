
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


## func DiffArray
``` go
func DiffArray(a, b []string) (diff []string)
```
DiffArray returns the delta of two arrays


## func GetHTTP
``` go
func GetHTTP(method, URL string) (io.Reader, error)
```
GetHTTP creates http requests to download data


## func LoadCfg
``` go
func LoadCfg() (string, error)
```
LoadCfg returns an EdgeOS config file string and error


## func ToBool
``` go
func ToBool(s string) bool
```
ToBool converts a string ("true" or "false") to it's boolean equivalent



## type CFile
``` go
type CFile struct {
    *Parms
    // contains filtered or unexported fields
}
```
CFile holds an array of file names











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
func ReadCfg(reader io.Reader) (*Config, error)
```
ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure




### func (\*Config) Get
``` go
func (c *Config) Get(node string) (o *Object)
```
Get returns an *Object for a given node



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



## type Content
``` go
type Content struct {
    *Object
    // contains filtered or unexported fields
}
```
Content is a struct of blacklist content











### func (Content) Source
``` go
func (d Content) Source(ltype string) *Objects
```
Source returns a map of sources



## type Contents
``` go
type Contents []*Content
```
Contents is an array of *content











## type Keys
``` go
type Keys []string
```
Keys is used for sorting operations on map Keys











### func (Keys) Len
``` go
func (k Keys) Len() int
```
len returns length of Keys



### func (Keys) Less
``` go
func (k Keys) Less(i, j int) bool
```
less returns the smallest element



### func (Keys) Swap
``` go
func (k Keys) Swap(i, j int)
```
Swap swaps elements of a key array



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
Includes returns a Content struct of blacklist Includes



### func (Object) Source
``` go
func (d Object) Source(ltype string) *Objects
```
Source returns a map of sources



### func (\*Object) String
``` go
func (o *Object) String() (r string)
```


## type Objects
``` go
type Objects struct {
    S []*Object
    *Parms
}
```
Objects is a struct of []*Objects











### func (Objects) Files
``` go
func (o Objects) Files() *CFile
```
Files returns a list of dnsmasq conf files from all srcs



### func (\*Objects) GetContent
``` go
func (objs *Objects) GetContent() *Contents
```
GetContent returns a Content struct



## type Option
``` go
type Option func(p *Parms) Option
```
Option sets is a recursive function









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
    // contains filtered or unexported fields
}
```
Parms is struct of parameters









### func NewParms
``` go
func NewParms(c *Config) *Parms
```
NewParms sets a new *Parms instance




### func (\*Parms) SetOpt
``` go
func (p *Parms) SetOpt(opts ...Option) (previous Option)
```
SetOpt sets the specified options passed as Parms and returns an option to restore the last arg's previous value



### func (Parms) String
``` go
func (p Parms) String() string
```
String method to implement fmt.Print interface









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)