
# edgeos
    import "github.com/britannic/blacklist/internal/edgeos"

Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.




## Constants
``` go
const (

    // ExcDomns labels domain exclusions
    ExcDomns = "domn-excludes"
    // ExcHosts labels host exclusions
    ExcHosts = "host-excludes"
    // ExcRoots labels global domain exclusions
    ExcRoots = "root-excludes"
    // False is a string constant
    False = "false"
    // PreDomns designates string label for preconfigured blacklisted domains
    PreDomns = preNoun + "-domain"
    // PreHosts designates string label for preconfigured blacklisted hosts
    PreHosts = preNoun + "-host"
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


## func StrToBool
``` go
func StrToBool(s string) bool
```
StrToBool converts a string ("true" or "false") to boolean



## type CFGcli
``` go
type CFGcli struct {
    *Config
    Cfg string
}
```
CFGcli loads configurations using the EdgeOS CFGcli











## type CFGstatic
``` go
type CFGstatic struct {
    *Config
    Cfg string
}
```
CFGstatic loads static configurations for testing











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



## type ConfLoader
``` go
type ConfLoader interface {
    // contains filtered or unexported methods
}
```
ConfLoader interface defines configuration load method











## type Config
``` go
type Config struct {
    *Parms
    // contains filtered or unexported fields
}
```
Config is a struct of configuration fields









### func NewConfig
``` go
func NewConfig(opts ...Option) *Config
```
NewConfig returns a new *Config initialized with the parameter options passed to it




### func (\*Config) Get
``` go
func (c *Config) Get(node string) *Objects
```
Get returns an *Object for a given node



### func (\*Config) GetAll
``` go
func (c *Config) GetAll(ltypes ...string) *Objects
```
GetAll returns an array of Objects



### func (\*Config) InSession
``` go
func (c *Config) InSession() bool
```
InSession returns true if VyOS/EdgeOS configuration is in session



### func (\*Config) LTypes
``` go
func (c *Config) LTypes() []string
```
LTypes returns an array of configured nodes



### func (\*Config) NewContent
``` go
func (c *Config) NewContent(iface IFace) (Contenter, error)
```
NewContent returns an interface of the requested IFace type



### func (\*Config) Nodes
``` go
func (c *Config) Nodes() (nodes []string)
```
Nodes returns an array of configured nodes



### func (\*Config) ProcessContent
``` go
func (c *Config) ProcessContent(cts ...Contenter) error
```
ProcessContent processes the Contents array



### func (\*Config) ReadCfg
``` go
func (c *Config) ReadCfg(r ConfLoader) error
```
ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure



### func (\*Config) ReloadDNS
``` go
func (c *Config) ReloadDNS() ([]byte, error)
```
ReloadDNS reloads the dnsmasq configuration



### func (\*Config) SetOpt
``` go
func (c *Config) SetOpt(opts ...Option) (previous Option)
```
SetOpt sets the specified options passed as Parms and returns an option to restore the last set of arg's previous values



### func (\*Config) String
``` go
func (c *Config) String() (s string)
```
String returns pretty print for the Blacklist struct



## type Contenter
``` go
type Contenter interface {
    Find(elem string) int
    GetList() *Objects
    SetURL(name string, url string)
    String() string
}
```
Contenter is a Content interface











## type ExcDomnObjects
``` go
type ExcDomnObjects struct {
    *Objects
}
```
ExcDomnObjects implements GetList for domain exclusions











### func (\*ExcDomnObjects) Find
``` go
func (e *ExcDomnObjects) Find(elem string) int
```
Find returns the int position of an Objects' element



### func (\*ExcDomnObjects) GetList
``` go
func (e *ExcDomnObjects) GetList() *Objects
```
GetList implements the Contenter interface for ExcDomnObjects



### func (\*ExcDomnObjects) SetURL
``` go
func (e *ExcDomnObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value



### func (\*ExcDomnObjects) String
``` go
func (e *ExcDomnObjects) String() string
```


## type ExcHostObjects
``` go
type ExcHostObjects struct {
    *Objects
}
```
ExcHostObjects implements GetList for host exclusions











### func (\*ExcHostObjects) Find
``` go
func (e *ExcHostObjects) Find(elem string) int
```
Find returns the int position of an Objects' element



### func (\*ExcHostObjects) GetList
``` go
func (e *ExcHostObjects) GetList() *Objects
```
GetList implements the Contenter interface for ExcHostObjects



### func (\*ExcHostObjects) SetURL
``` go
func (e *ExcHostObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value



### func (\*ExcHostObjects) String
``` go
func (e *ExcHostObjects) String() string
```


## type ExcRootObjects
``` go
type ExcRootObjects struct {
    *Objects
}
```
ExcRootObjects implements GetList for global domain exclusions











### func (\*ExcRootObjects) Find
``` go
func (e *ExcRootObjects) Find(elem string) int
```
Find returns the int position of an Objects' element



### func (\*ExcRootObjects) GetList
``` go
func (e *ExcRootObjects) GetList() *Objects
```
GetList implements the Contenter interface for ExcRootObjects



### func (\*ExcRootObjects) SetURL
``` go
func (e *ExcRootObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value



### func (\*ExcRootObjects) String
``` go
func (e *ExcRootObjects) String() string
```


## type FIODataObjects
``` go
type FIODataObjects struct {
    *Objects
}
```
FIODataObjects implements GetList for files











### func (\*FIODataObjects) Find
``` go
func (f *FIODataObjects) Find(elem string) int
```
Find returns the int position of an Objects' element



### func (\*FIODataObjects) GetList
``` go
func (f *FIODataObjects) GetList() *Objects
```
GetList implements the Contenter interface for FIODataObjects



### func (\*FIODataObjects) SetURL
``` go
func (f *FIODataObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value



### func (\*FIODataObjects) String
``` go
func (f *FIODataObjects) String() string
```


## type IFace
``` go
type IFace int
```
IFace type for labeling interface types



``` go
const (
    Invalid IFace = iota + 100
    ExRtObj
    ExDmObj
    ExHtObj
    FileObj
    PreDObj
    PreHObj
    URLsObj
)
```
IFace types for labeling interface types









### func (IFace) String
``` go
func (i IFace) String() (s string)
```


## type Objects
``` go
type Objects struct {
    *Parms
    // contains filtered or unexported fields
}
```
Objects is a struct of []*Object











### func (\*Objects) Files
``` go
func (o *Objects) Files() *CFile
```
Files returns a list of dnsmasq conf files from all srcs



### func (\*Objects) Find
``` go
func (o *Objects) Find(elem string) int
```
Find returns the int position of an Objects' element



### func (\*Objects) Len
``` go
func (o *Objects) Len() int
```
Implement Sort Interface for Objects



### func (\*Objects) Less
``` go
func (o *Objects) Less(i, j int) bool
```


### func (\*Objects) Names
``` go
func (o *Objects) Names() (s sort.StringSlice)
```
Names returns a sorted slice of Objects names



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
type Option func(c *Config) Option
```
Option is a recursive function









### func API
``` go
func API(s string) Option
```
API sets the EdgeOS CLI API command


### func Arch
``` go
func Arch(arch string) Option
```
Arch sets target CPU architecture


### func Bash
``` go
func Bash(cmd string) Option
```
Bash sets the shell processor


### func Cores
``` go
func Cores(i int) Option
```
Cores sets max CPU cores


### func DNSsvc
``` go
func DNSsvc(d string) Option
```
DNSsvc sets dnsmasq restart command


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


### func InCLI
``` go
func InCLI(in string) Option
```
InCLI sets the CLI inSession command


### func LTypes
``` go
func LTypes(s []string) Option
```
LTypes sets an array of legal types used by Source


### func Level
``` go
func Level(s string) Option
```
Level sets the EdgeOS API CLI level


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


### func Test
``` go
func Test(b bool) Option
```
Test toggles testing mode on or off


### func Timeout
``` go
func Timeout(t time.Duration) Option
```
Timeout sets how long before an unresponsive goroutine is aborted


### func Verbosity
``` go
func Verbosity(i int) Option
```
Verbosity sets the verbosity level to v


### func WCard
``` go
func WCard(w Wildcard) Option
```
WCard sets file globbing wildcard values




## type Parms
``` go
type Parms struct {
    Wildcard
    API       string
    Arch      string
    Bash      string
    Cores     int
    Debug     bool
    Dex       list
    Dir       string
    DNSsvc    string
    Exc       list
    Ext       string
    File      string
    FnFmt     string
    InCLI     string
    Level     string
    Method    string
    Nodes     []string
    Pfx       string
    Poll      int
    Ltypes    []string
    Test      bool
    Timeout   time.Duration
    Verbosity int
}
```
Parms is struct of parameters









### func NewParms
``` go
func NewParms() *Parms
```
NewParms sets a new *Parms instance




### func (\*Parms) String
``` go
func (p *Parms) String() string
```
String method to implement fmt.Print interface



## type PreDomnObjects
``` go
type PreDomnObjects struct {
    *Objects
}
```
PreDomnObjects implements GetList for pre-configured domains content











### func (\*PreDomnObjects) Find
``` go
func (p *PreDomnObjects) Find(elem string) int
```
Find returns the int position of an Objects' element



### func (\*PreDomnObjects) GetList
``` go
func (p *PreDomnObjects) GetList() *Objects
```
GetList implements the Contenter interface for PreDomnObjects



### func (\*PreDomnObjects) SetURL
``` go
func (p *PreDomnObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value



### func (\*PreDomnObjects) String
``` go
func (p *PreDomnObjects) String() string
```


## type PreHostObjects
``` go
type PreHostObjects struct {
    *Objects
}
```
PreHostObjects implements GetList for pre-configured hosts content











### func (\*PreHostObjects) Find
``` go
func (p *PreHostObjects) Find(elem string) int
```
Find returns the int position of an Objects' element



### func (\*PreHostObjects) GetList
``` go
func (p *PreHostObjects) GetList() *Objects
```
GetList implements the Contenter interface for PreHostObjects



### func (\*PreHostObjects) SetURL
``` go
func (p *PreHostObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value



### func (\*PreHostObjects) String
``` go
func (p *PreHostObjects) String() string
```


## type URLDataObjects
``` go
type URLDataObjects struct {
    *Objects
}
```
URLDataObjects implements GetList for URLs











### func (\*URLDataObjects) Find
``` go
func (u *URLDataObjects) Find(elem string) int
```
Find returns the int position of an Objects' element



### func (\*URLDataObjects) GetList
``` go
func (u *URLDataObjects) GetList() *Objects
```
GetList implements the Contenter interface for URLDataObjects



### func (\*URLDataObjects) SetURL
``` go
func (u *URLDataObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value



### func (\*URLDataObjects) String
``` go
func (u *URLDataObjects) String() string
```


## type Wildcard
``` go
type Wildcard struct {
    Node string
    Name string
}
```
Wildcard struct sets globbing wildcards for filename searches

















- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)