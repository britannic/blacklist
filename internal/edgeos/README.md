

# edgeos
`import "github.com/britannic/blacklist/internal/edgeos"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.

Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.

Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [func GetFile(f string) (io.Reader, error)](#GetFile)
* [func Iter(i int) []struct{}](#Iter)
* [func NewWriter() io.Writer](#NewWriter)
* [type CFGcli](#CFGcli)
* [type CFGstatic](#CFGstatic)
* [type CFile](#CFile)
  * [func (c *CFile) Remove() error](#CFile.Remove)
  * [func (c *CFile) String() string](#CFile.String)
  * [func (c *CFile) Strings() []string](#CFile.Strings)
* [type ConfLoader](#ConfLoader)
* [type Config](#Config)
  * [func NewConfig(opts ...Option) *Config](#NewConfig)
  * [func (c *Config) Get(node string) *Objects](#Config.Get)
  * [func (c *Config) GetAll(ltypes ...string) *Objects](#Config.GetAll)
  * [func (c *Config) GetTotalStats() (dropped, extracted, kept int32)](#Config.GetTotalStats)
  * [func (c *Config) InSession() bool](#Config.InSession)
  * [func (c *Config) NewContent(iface IFace) (Contenter, error)](#Config.NewContent)
  * [func (c *Config) Nodes() (nodes []string)](#Config.Nodes)
  * [func (c *Config) ProcessContent(cts ...Contenter) error](#Config.ProcessContent)
  * [func (c *Config) ReadCfg(r ConfLoader) error](#Config.ReadCfg)
  * [func (c *Config) ReloadDNS() ([]byte, error)](#Config.ReloadDNS)
  * [func (c *Config) SetOpt(opts ...Option) Option](#Config.SetOpt)
  * [func (c *Config) String() (s string)](#Config.String)
* [type Contenter](#Contenter)
* [type Env](#Env)
  * [func (e *Env) Debug(s ...interface{})](#Env.Debug)
  * [func (e *Env) String() string](#Env.String)
* [type ExcDomnObjects](#ExcDomnObjects)
  * [func (e *ExcDomnObjects) Find(s string) int](#ExcDomnObjects.Find)
  * [func (e *ExcDomnObjects) GetList() *Objects](#ExcDomnObjects.GetList)
  * [func (e *ExcDomnObjects) Len() int](#ExcDomnObjects.Len)
  * [func (e *ExcDomnObjects) SetURL(name, url string)](#ExcDomnObjects.SetURL)
  * [func (e *ExcDomnObjects) String() string](#ExcDomnObjects.String)
* [type ExcHostObjects](#ExcHostObjects)
  * [func (e *ExcHostObjects) Find(s string) int](#ExcHostObjects.Find)
  * [func (e *ExcHostObjects) GetList() *Objects](#ExcHostObjects.GetList)
  * [func (e *ExcHostObjects) Len() int](#ExcHostObjects.Len)
  * [func (e *ExcHostObjects) SetURL(name, url string)](#ExcHostObjects.SetURL)
  * [func (e *ExcHostObjects) String() string](#ExcHostObjects.String)
* [type ExcRootObjects](#ExcRootObjects)
  * [func (e *ExcRootObjects) Find(s string) int](#ExcRootObjects.Find)
  * [func (e *ExcRootObjects) GetList() *Objects](#ExcRootObjects.GetList)
  * [func (e *ExcRootObjects) Len() int](#ExcRootObjects.Len)
  * [func (e *ExcRootObjects) SetURL(name, url string)](#ExcRootObjects.SetURL)
  * [func (e *ExcRootObjects) String() string](#ExcRootObjects.String)
* [type FIODataObjects](#FIODataObjects)
  * [func (f *FIODataObjects) Find(s string) int](#FIODataObjects.Find)
  * [func (f *FIODataObjects) GetList() *Objects](#FIODataObjects.GetList)
  * [func (f *FIODataObjects) Len() int](#FIODataObjects.Len)
  * [func (f *FIODataObjects) SetURL(name, url string)](#FIODataObjects.SetURL)
  * [func (f *FIODataObjects) String() string](#FIODataObjects.String)
* [type IFace](#IFace)
  * [func (i IFace) String() string](#IFace.String)
* [type Objects](#Objects)
  * [func (o *Objects) Files() *CFile](#Objects.Files)
  * [func (o *Objects) Filter(ltype string) *Objects](#Objects.Filter)
  * [func (o *Objects) Find(elem string) int](#Objects.Find)
  * [func (o *Objects) Len() int](#Objects.Len)
  * [func (o *Objects) Less(i, j int) bool](#Objects.Less)
  * [func (o *Objects) Names() (s sort.StringSlice)](#Objects.Names)
  * [func (o *Objects) String() (s string)](#Objects.String)
  * [func (o *Objects) Swap(i, j int)](#Objects.Swap)
* [type Option](#Option)
  * [func API(s string) Option](#API)
  * [func Arch(s string) Option](#Arch)
  * [func Bash(s string) Option](#Bash)
  * [func Cores(i int) Option](#Cores)
  * [func DNSsvc(s string) Option](#DNSsvc)
  * [func Dbug(b bool) Option](#Dbug)
  * [func Dir(s string) Option](#Dir)
  * [func Disabled(b bool) Option](#Disabled)
  * [func Ext(s string) Option](#Ext)
  * [func File(s string) Option](#File)
  * [func FileNameFmt(s string) Option](#FileNameFmt)
  * [func InCLI(s string) Option](#InCLI)
  * [func Level(s string) Option](#Level)
  * [func Logger(l *logging.Logger) Option](#Logger)
  * [func Method(s string) Option](#Method)
  * [func Prefix(d string, h string) Option](#Prefix)
  * [func Test(b bool) Option](#Test)
  * [func Timeout(t time.Duration) Option](#Timeout)
  * [func Verb(b bool) Option](#Verb)
  * [func WCard(w Wildcard) Option](#WCard)
  * [func Writer(w io.Writer) Option](#Writer)
* [type PreDomnObjects](#PreDomnObjects)
  * [func (p *PreDomnObjects) Find(s string) int](#PreDomnObjects.Find)
  * [func (p *PreDomnObjects) GetList() *Objects](#PreDomnObjects.GetList)
  * [func (p *PreDomnObjects) Len() int](#PreDomnObjects.Len)
  * [func (p *PreDomnObjects) SetURL(name, url string)](#PreDomnObjects.SetURL)
  * [func (p *PreDomnObjects) String() string](#PreDomnObjects.String)
* [type PreHostObjects](#PreHostObjects)
  * [func (p *PreHostObjects) Find(s string) int](#PreHostObjects.Find)
  * [func (p *PreHostObjects) GetList() *Objects](#PreHostObjects.GetList)
  * [func (p *PreHostObjects) Len() int](#PreHostObjects.Len)
  * [func (p *PreHostObjects) SetURL(name, url string)](#PreHostObjects.SetURL)
  * [func (p *PreHostObjects) String() string](#PreHostObjects.String)
* [type PreRootObjects](#PreRootObjects)
  * [func (p *PreRootObjects) Find(s string) int](#PreRootObjects.Find)
  * [func (p *PreRootObjects) GetList() *Objects](#PreRootObjects.GetList)
  * [func (p *PreRootObjects) Len() int](#PreRootObjects.Len)
  * [func (p *PreRootObjects) SetURL(name, url string)](#PreRootObjects.SetURL)
  * [func (p *PreRootObjects) String() string](#PreRootObjects.String)
* [type URLDomnObjects](#URLDomnObjects)
  * [func (u *URLDomnObjects) Find(s string) int](#URLDomnObjects.Find)
  * [func (u *URLDomnObjects) GetList() *Objects](#URLDomnObjects.GetList)
  * [func (u *URLDomnObjects) Len() int](#URLDomnObjects.Len)
  * [func (u *URLDomnObjects) SetURL(name, url string)](#URLDomnObjects.SetURL)
  * [func (u *URLDomnObjects) String() string](#URLDomnObjects.String)
* [type URLHostObjects](#URLHostObjects)
  * [func (u *URLHostObjects) Find(s string) int](#URLHostObjects.Find)
  * [func (u *URLHostObjects) GetList() *Objects](#URLHostObjects.GetList)
  * [func (u *URLHostObjects) Len() int](#URLHostObjects.Len)
  * [func (u *URLHostObjects) SetURL(name, url string)](#URLHostObjects.SetURL)
  * [func (u *URLHostObjects) String() string](#URLHostObjects.String)
* [type Wildcard](#Wildcard)


#### <a name="pkg-files">Package files</a>
[cfile.go](/src/github.com/britannic/blacklist/internal/edgeos/cfile.go) [config.go](/src/github.com/britannic/blacklist/internal/edgeos/config.go) [content.go](/src/github.com/britannic/blacklist/internal/edgeos/content.go) [data.go](/src/github.com/britannic/blacklist/internal/edgeos/data.go) [http.go](/src/github.com/britannic/blacklist/internal/edgeos/http.go) [io.go](/src/github.com/britannic/blacklist/internal/edgeos/io.go) [json.go](/src/github.com/britannic/blacklist/internal/edgeos/json.go) [list.go](/src/github.com/britannic/blacklist/internal/edgeos/list.go) [ntype_string.go](/src/github.com/britannic/blacklist/internal/edgeos/ntype_string.go) [object.go](/src/github.com/britannic/blacklist/internal/edgeos/object.go) [opts.go](/src/github.com/britannic/blacklist/internal/edgeos/opts.go) [source.go](/src/github.com/britannic/blacklist/internal/edgeos/source.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (

    // ExcDomns is a string labels for domain exclusions
    ExcDomns = "whitelisted-subdomains"
    // ExcHosts is a string labels for host exclusions
    ExcHosts = "whitelisted-servers"
    // ExcRoots is a string labels for preconfigured global domain exclusions
    ExcRoots = "global-whitelisted-domains"
    // PreDomns is a string label for preconfigured whitelisted domains
    PreDomns = "blacklisted-subdomains"
    // PreHosts is a string label for preconfigured blacklisted hosts
    PreHosts = "blacklisted-servers"
    // PreRoots is a string label for preconfigured global blacklisted hosts
    PreRoots = "global-blacklisted-domains"
    // False is a string constant
    False = "false"
    // True is a string constant
    True = "true"
)
```



## <a name="GetFile">func</a> [GetFile](/src/target/io.go?s=1940:1981#L93)
``` go
func GetFile(f string) (io.Reader, error)
```
GetFile reads a file and returns an io.Reader



## <a name="Iter">func</a> [Iter](/src/target/data.go?s=2299:2326#L102)
``` go
func Iter(i int) []struct{}
```
Iter iterates over ints - use it in for loops



## <a name="NewWriter">func</a> [NewWriter](/src/target/data.go?s=2394:2420#L107)
``` go
func NewWriter() io.Writer
```
NewWriter returns an io.Writer




## <a name="CFGcli">type</a> [CFGcli](/src/target/io.go?s=139:182#L14)
``` go
type CFGcli struct {
    *Config
    Cfg string
}
```
CFGcli loads configurations using the EdgeOS CFGcli










## <a name="CFGstatic">type</a> [CFGstatic](/src/target/io.go?s=237:283#L20)
``` go
type CFGstatic struct {
    *Config
    Cfg string
}
```
CFGstatic loads static configurations for testing










## <a name="CFile">type</a> [CFile](/src/target/cfile.go?s=226:282#L12)
``` go
type CFile struct {
    *Env
    Names []string
    // contains filtered or unexported fields
}
```
CFile holds an array of file names










### <a name="CFile.Remove">func</a> (\*CFile) [Remove](/src/target/cfile.go?s=577:607#L26)
``` go
func (c *CFile) Remove() error
```
Remove deletes a CFile array of file names




### <a name="CFile.String">func</a> (\*CFile) [String](/src/target/cfile.go?s=875:906#L37)
``` go
func (c *CFile) String() string
```
String implements string method




### <a name="CFile.Strings">func</a> (\*CFile) [Strings](/src/target/cfile.go?s=998:1032#L42)
``` go
func (c *CFile) Strings() []string
```
Strings returns a sorted array of strings.




## <a name="ConfLoader">type</a> [ConfLoader](/src/target/config.go?s=407:454#L23)
``` go
type ConfLoader interface {
    // contains filtered or unexported methods
}
```
ConfLoader interface handles multiple configuration load methods










## <a name="Config">type</a> [Config](/src/target/config.go?s=502:536#L28)
``` go
type Config struct {
    *Env
    // contains filtered or unexported fields
}
```
Config is a struct of configuration fields







### <a name="NewConfig">func</a> [NewConfig](/src/target/opts.go?s=4833:4871#L216)
``` go
func NewConfig(opts ...Option) *Config
```
NewConfig returns a new *Config initialized with the parameter options passed to it





### <a name="Config.Get">func</a> (\*Config) [Get](/src/target/config.go?s=4845:4887#L219)
``` go
func (c *Config) Get(node string) *Objects
```
Get returns an *Object for a given node




### <a name="Config.GetAll">func</a> (\*Config) [GetAll](/src/target/config.go?s=5108:5158#L233)
``` go
func (c *Config) GetAll(ltypes ...string) *Objects
```
GetAll returns an array of Objects




### <a name="Config.GetTotalStats">func</a> (\*Config) [GetTotalStats](/src/target/config.go?s=3024:3089#L154)
``` go
func (c *Config) GetTotalStats() (dropped, extracted, kept int32)
```
GetTotalStats displays aggregate statistics for processed sources




### <a name="Config.InSession">func</a> (\*Config) [InSession](/src/target/config.go?s=5338:5371#L242)
``` go
func (c *Config) InSession() bool
```
InSession returns true if VyOS/EdgeOS configure is in session




### <a name="Config.NewContent">func</a> (\*Config) [NewContent](/src/target/config.go?s=3550:3609#L173)
``` go
func (c *Config) NewContent(iface IFace) (Contenter, error)
```
NewContent returns an interface of the requested IFace type




### <a name="Config.Nodes">func</a> (\*Config) [Nodes](/src/target/config.go?s=5896:5937#L263)
``` go
func (c *Config) Nodes() (nodes []string)
```
Nodes returns an array of configured nodes




### <a name="Config.ProcessContent">func</a> (\*Config) [ProcessContent](/src/target/config.go?s=7914:7969#L346)
``` go
func (c *Config) ProcessContent(cts ...Contenter) error
```
ProcessContent processes the Contents array




### <a name="Config.ReadCfg">func</a> (\*Config) [ReadCfg](/src/target/config.go?s=9072:9116#L402)
``` go
func (c *Config) ReadCfg(r ConfLoader) error
```
ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure




### <a name="Config.ReloadDNS">func</a> (\*Config) [ReloadDNS](/src/target/config.go?s=11091:11135#L459)
``` go
func (c *Config) ReloadDNS() ([]byte, error)
```
ReloadDNS reloads the dnsmasq configuration




### <a name="Config.SetOpt">func</a> (\*Config) [SetOpt](/src/target/opts.go?s=1885:1931#L64)
``` go
func (c *Config) SetOpt(opts ...Option) Option
```
SetOpt sets the specified options passed as Env and returns an option to restore the last set of arg's previous values




### <a name="Config.String">func</a> (\*Config) [String](/src/target/config.go?s=11553:11589#L478)
``` go
func (c *Config) String() (s string)
```
String returns pretty print for the Blacklist struct




## <a name="Contenter">type</a> [Contenter](/src/target/content.go?s=452:570#L37)
``` go
type Contenter interface {
    Find(string) int
    GetList() *Objects
    Len() int
    SetURL(string, string)
    String() string
}
```
Contenter is a Content interface










## <a name="Env">type</a> [Env](/src/target/opts.go?s=153:1266#L14)
``` go
type Env struct {
    Log      *logging.Logger
    API      string        `json:"API,omitempty"`
    Arch     string        `json:"Arch,omitempty"`
    Bash     string        `json:"Bash,omitempty"`
    Cores    int           `json:"Cores,omitempty"`
    Disabled bool          `json:"Disabled"`
    Dbug     bool          `json:"Dbug,omitempty"`
    Dex      list          `json:"Dex,omitempty"`
    Dir      string        `json:"Dir,omitempty"`
    DNSsvc   string        `json:"dnsmasq service,omitempty"`
    Exc      list          `json:"Exc,omitempty"`
    Ext      string        `json:"dnsmasq fileExt.,omitempty"`
    File     string        `json:"File,omitempty"`
    FnFmt    string        `json:"File name fmt,omitempty"`
    InCLI    string        `json:"-"`
    Level    string        `json:"CLI Path,omitempty"`
    Method   string        `json:"HTTP method,omitempty"`
    Pfx      dnsPfx        `json:"Prefix,omitempty"`
    Test     bool          `json:"Test,omitempty"`
    Timeout  time.Duration `json:"Timeout,omitempty"`
    Verb     bool          `json:"Verbosity,omitempty"`
    Wildcard `json:"Wildcard,omitempty"`
    // contains filtered or unexported fields
}
```
Env is struct of parameters










### <a name="Env.Debug">func</a> (\*Env) [Debug](/src/target/opts.go?s=1614:1651#L54)
``` go
func (e *Env) Debug(s ...interface{})
```
Debug logs debug messages when the Dbug flag is true




### <a name="Env.String">func</a> (\*Env) [String](/src/target/opts.go?s=5429:5458#L242)
``` go
func (e *Env) String() string
```
Env Stringer interface




## <a name="ExcDomnObjects">type</a> [ExcDomnObjects](/src/target/content.go?s=631:671#L46)
``` go
type ExcDomnObjects struct {
    *Objects
}
```
ExcDomnObjects struct of *Objects for domain exclusions










### <a name="ExcDomnObjects.Find">func</a> (\*ExcDomnObjects) [Find](/src/target/content.go?s=1551:1594#L91)
``` go
func (e *ExcDomnObjects) Find(s string) int
```
Find returns the int position of an Objects' element




### <a name="ExcDomnObjects.GetList">func</a> (\*ExcDomnObjects) [GetList](/src/target/content.go?s=3243:3286#L181)
``` go
func (e *ExcDomnObjects) GetList() *Objects
```
GetList implements the Contenter interface for ExcDomnObjects




### <a name="ExcDomnObjects.Len">func</a> (\*ExcDomnObjects) [Len](/src/target/content.go?s=5946:5980#L312)
``` go
func (e *ExcDomnObjects) Len() int
```
Len returns how many sources there are




### <a name="ExcDomnObjects.SetURL">func</a> (\*ExcDomnObjects) [SetURL](/src/target/content.go?s=8086:8135#L390)
``` go
func (e *ExcDomnObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="ExcDomnObjects.String">func</a> (\*ExcDomnObjects) [String](/src/target/content.go?s=9572:9612#L470)
``` go
func (e *ExcDomnObjects) String() string
```



## <a name="ExcHostObjects">type</a> [ExcHostObjects](/src/target/content.go?s=730:770#L51)
``` go
type ExcHostObjects struct {
    *Objects
}
```
ExcHostObjects struct of *Objects for host exclusions










### <a name="ExcHostObjects.Find">func</a> (\*ExcHostObjects) [Find](/src/target/content.go?s=1738:1781#L101)
``` go
func (e *ExcHostObjects) Find(s string) int
```
Find returns the int position of an Objects' element




### <a name="ExcHostObjects.GetList">func</a> (\*ExcHostObjects) [GetList](/src/target/content.go?s=3502:3545#L194)
``` go
func (e *ExcHostObjects) GetList() *Objects
```
GetList implements the Contenter interface for ExcHostObjects




### <a name="ExcHostObjects.Len">func</a> (\*ExcHostObjects) [Len](/src/target/content.go?s=6046:6080#L315)
``` go
func (e *ExcHostObjects) Len() int
```
Len returns how many sources there are




### <a name="ExcHostObjects.SetURL">func</a> (\*ExcHostObjects) [SetURL](/src/target/content.go?s=8256:8305#L399)
``` go
func (e *ExcHostObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="ExcHostObjects.String">func</a> (\*ExcHostObjects) [String](/src/target/content.go?s=9643:9683#L471)
``` go
func (e *ExcHostObjects) String() string
```



## <a name="ExcRootObjects">type</a> [ExcRootObjects](/src/target/content.go?s=838:878#L56)
``` go
type ExcRootObjects struct {
    *Objects
}
```
ExcRootObjects struct of *Objects for global domain exclusions










### <a name="ExcRootObjects.Find">func</a> (\*ExcRootObjects) [Find](/src/target/content.go?s=1925:1968#L111)
``` go
func (e *ExcRootObjects) Find(s string) int
```
Find returns the int position of an Objects' element




### <a name="ExcRootObjects.GetList">func</a> (\*ExcRootObjects) [GetList](/src/target/content.go?s=3761:3804#L207)
``` go
func (e *ExcRootObjects) GetList() *Objects
```
GetList implements the Contenter interface for ExcRootObjects




### <a name="ExcRootObjects.Len">func</a> (\*ExcRootObjects) [Len](/src/target/content.go?s=6146:6180#L318)
``` go
func (e *ExcRootObjects) Len() int
```
Len returns how many sources there are




### <a name="ExcRootObjects.SetURL">func</a> (\*ExcRootObjects) [SetURL](/src/target/content.go?s=8426:8475#L408)
``` go
func (e *ExcRootObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="ExcRootObjects.String">func</a> (\*ExcRootObjects) [String](/src/target/content.go?s=9714:9754#L472)
``` go
func (e *ExcRootObjects) String() string
```



## <a name="FIODataObjects">type</a> [FIODataObjects](/src/target/content.go?s=927:967#L61)
``` go
type FIODataObjects struct {
    *Objects
}
```
FIODataObjects struct of *Objects for files










### <a name="FIODataObjects.Find">func</a> (\*FIODataObjects) [Find](/src/target/content.go?s=2112:2155#L121)
``` go
func (f *FIODataObjects) Find(s string) int
```
Find returns the int position of an Objects' element




### <a name="FIODataObjects.GetList">func</a> (\*FIODataObjects) [GetList](/src/target/content.go?s=4020:4063#L220)
``` go
func (f *FIODataObjects) GetList() *Objects
```
GetList implements the Contenter interface for FIODataObjects




### <a name="FIODataObjects.Len">func</a> (\*FIODataObjects) [Len](/src/target/content.go?s=6246:6280#L321)
``` go
func (f *FIODataObjects) Len() int
```
Len returns how many sources there are




### <a name="FIODataObjects.SetURL">func</a> (\*FIODataObjects) [SetURL](/src/target/content.go?s=8596:8645#L417)
``` go
func (f *FIODataObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="FIODataObjects.String">func</a> (\*FIODataObjects) [String](/src/target/content.go?s=9785:9825#L473)
``` go
func (f *FIODataObjects) String() string
```



## <a name="IFace">type</a> [IFace](/src/target/content.go?s=152:166#L13)
``` go
type IFace int
```
IFace type for labeling interface types


``` go
const (
    Invalid IFace = iota + 100
    ExDmObj
    ExHtObj
    ExRtObj
    FileObj
    PreDObj
    PreHObj
    PreRObj
    URLdObj
    URLhObj
)
```
IFace types for labeling interface types










### <a name="IFace.String">func</a> (IFace) [String](/src/target/content.go?s=10212:10242#L480)
``` go
func (i IFace) String() string
```



## <a name="Objects">type</a> [Objects](/src/target/object.go?s=90:134#L10)
``` go
type Objects struct {
    *Env
    // contains filtered or unexported fields
}
```
Objects is a struct of []*source










### <a name="Objects.Files">func</a> (\*Objects) [Files](/src/target/object.go?s=342:374#L21)
``` go
func (o *Objects) Files() *CFile
```
Files returns a list of dnsmasq conf files from all srcs




### <a name="Objects.Filter">func</a> (\*Objects) [Filter](/src/target/object.go?s=655:702#L34)
``` go
func (o *Objects) Filter(ltype string) *Objects
```
Filter returns a subset of Objects filtered by ltype




### <a name="Objects.Find">func</a> (\*Objects) [Find](/src/target/object.go?s=1109:1148#L54)
``` go
func (o *Objects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="Objects.Len">func</a> (\*Objects) [Len](/src/target/object.go?s=2776:2803#L138)
``` go
func (o *Objects) Len() int
```
Implement Sort Interface for Objects




### <a name="Objects.Less">func</a> (\*Objects) [Less](/src/target/object.go?s=2836:2873#L139)
``` go
func (o *Objects) Less(i, j int) bool
```



### <a name="Objects.Names">func</a> (\*Objects) [Names](/src/target/object.go?s=2473:2519#L121)
``` go
func (o *Objects) Names() (s sort.StringSlice)
```
Names returns a sorted slice of Objects names




### <a name="Objects.String">func</a> (\*Objects) [String](/src/target/object.go?s=2631:2668#L130)
``` go
func (o *Objects) String() (s string)
```
Stringer for Objects




### <a name="Objects.Swap">func</a> (\*Objects) [Swap](/src/target/object.go?s=2915:2947#L140)
``` go
func (o *Objects) Swap(i, j int)
```



## <a name="Option">type</a> [Option](/src/target/opts.go?s=1727:1761#L61)
``` go
type Option func(c *Config) Option
```
Option is a recursive function







### <a name="API">func</a> [API](/src/target/opts.go?s=2411:2436#L89)
``` go
func API(s string) Option
```
API sets the EdgeOS CLI API command


### <a name="Arch">func</a> [Arch](/src/target/opts.go?s=2246:2272#L80)
``` go
func Arch(s string) Option
```
Arch sets target CPU architecture


### <a name="Bash">func</a> [Bash](/src/target/opts.go?s=2566:2592#L98)
``` go
func Bash(s string) Option
```
Bash sets the shell processor


### <a name="Cores">func</a> [Cores](/src/target/opts.go?s=2720:2744#L107)
``` go
func Cores(i int) Option
```
Cores sets max CPU cores


### <a name="DNSsvc">func</a> [DNSsvc](/src/target/opts.go?s=3394:3422#L144)
``` go
func DNSsvc(s string) Option
```
DNSsvc sets dnsmasq restart command


### <a name="Dbug">func</a> [Dbug](/src/target/opts.go?s=3078:3102#L126)
``` go
func Dbug(b bool) Option
```
Dbug toggles Debug level on or off


### <a name="Dir">func</a> [Dir](/src/target/opts.go?s=3233:3258#L135)
``` go
func Dir(s string) Option
```
Dir sets directory location


### <a name="Disabled">func</a> [Disabled](/src/target/opts.go?s=2900:2928#L117)
``` go
func Disabled(b bool) Option
```
Disabled toggles Disabled


### <a name="Ext">func</a> [Ext](/src/target/opts.go?s=3571:3596#L153)
``` go
func Ext(s string) Option
```
Ext sets the blacklist file n extension


### <a name="File">func</a> [File](/src/target/opts.go?s=3736:3762#L162)
``` go
func File(s string) Option
```
File sets the EdgeOS configuration file


### <a name="FileNameFmt">func</a> [FileNameFmt](/src/target/opts.go?s=3924:3957#L171)
``` go
func FileNameFmt(s string) Option
```
FileNameFmt sets the EdgeOS configuration file name format


### <a name="InCLI">func</a> [InCLI](/src/target/opts.go?s=4106:4133#L180)
``` go
func InCLI(s string) Option
```
InCLI sets the CLI inSession command


### <a name="Level">func</a> [Level](/src/target/opts.go?s=4275:4302#L189)
``` go
func Level(s string) Option
```
Level sets the EdgeOS API CLI level


### <a name="Logger">func</a> [Logger](/src/target/opts.go?s=4444:4481#L198)
``` go
func Logger(l *logging.Logger) Option
```
Logger sets a pointer to the logger


### <a name="Method">func</a> [Method](/src/target/opts.go?s=4612:4640#L207)
``` go
func Method(s string) Option
```
Method sets the HTTP method


### <a name="Prefix">func</a> [Prefix](/src/target/opts.go?s=5222:5260#L232)
``` go
func Prefix(d string, h string) Option
```
Prefix sets the dnsmasq configuration address line prefix


### <a name="Test">func</a> [Test](/src/target/opts.go?s=5566:5590#L248)
``` go
func Test(b bool) Option
```
Test toggles testing mode on or off


### <a name="Timeout">func</a> [Timeout](/src/target/opts.go?s=5759:5795#L257)
``` go
func Timeout(t time.Duration) Option
```
Timeout sets how long before an unresponsive goroutine is aborted


### <a name="Verb">func</a> [Verb](/src/target/opts.go?s=5942:5966#L266)
``` go
func Verb(b bool) Option
```
Verb sets the verbosity level to v


### <a name="WCard">func</a> [WCard](/src/target/opts.go?s=6110:6139#L275)
``` go
func WCard(w Wildcard) Option
```
WCard sets file globbing wildcard values


### <a name="Writer">func</a> [Writer](/src/target/opts.go?s=6314:6345#L284)
``` go
func Writer(w io.Writer) Option
```
Writer provides an address for anything that can use io.Writer





## <a name="PreDomnObjects">type</a> [PreDomnObjects](/src/target/content.go?s=1041:1081#L66)
``` go
type PreDomnObjects struct {
    *Objects
}
```
PreDomnObjects struct of *Objects for pre-configured domains content










### <a name="PreDomnObjects.Find">func</a> (\*PreDomnObjects) [Find](/src/target/content.go?s=2299:2342#L131)
``` go
func (p *PreDomnObjects) Find(s string) int
```
Find returns the int position of an Objects' element




### <a name="PreDomnObjects.GetList">func</a> (\*PreDomnObjects) [GetList](/src/target/content.go?s=4436:4479#L240)
``` go
func (p *PreDomnObjects) GetList() *Objects
```
GetList implements the Contenter interface for PreDomnObjects




### <a name="PreDomnObjects.Len">func</a> (\*PreDomnObjects) [Len](/src/target/content.go?s=6346:6380#L324)
``` go
func (p *PreDomnObjects) Len() int
```
Len returns how many sources there are




### <a name="PreDomnObjects.SetURL">func</a> (\*PreDomnObjects) [SetURL](/src/target/content.go?s=8766:8815#L426)
``` go
func (p *PreDomnObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="PreDomnObjects.String">func</a> (\*PreDomnObjects) [String](/src/target/content.go?s=9856:9896#L474)
``` go
func (p *PreDomnObjects) String() string
```



## <a name="PreHostObjects">type</a> [PreHostObjects](/src/target/content.go?s=1153:1193#L71)
``` go
type PreHostObjects struct {
    *Objects
}
```
PreHostObjects struct of *Objects for pre-configured hosts content










### <a name="PreHostObjects.Find">func</a> (\*PreHostObjects) [Find](/src/target/content.go?s=2486:2529#L141)
``` go
func (p *PreHostObjects) Find(s string) int
```
Find returns the int position of an Objects' element




### <a name="PreHostObjects.GetList">func</a> (\*PreHostObjects) [GetList](/src/target/content.go?s=4684:4727#L251)
``` go
func (p *PreHostObjects) GetList() *Objects
```
GetList implements the Contenter interface for PreHostObjects




### <a name="PreHostObjects.Len">func</a> (\*PreHostObjects) [Len](/src/target/content.go?s=6446:6480#L327)
``` go
func (p *PreHostObjects) Len() int
```
Len returns how many sources there are




### <a name="PreHostObjects.SetURL">func</a> (\*PreHostObjects) [SetURL](/src/target/content.go?s=8936:8985#L435)
``` go
func (p *PreHostObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="PreHostObjects.String">func</a> (\*PreHostObjects) [String](/src/target/content.go?s=9927:9967#L475)
``` go
func (p *PreHostObjects) String() string
```



## <a name="PreRootObjects">type</a> [PreRootObjects](/src/target/content.go?s=1265:1305#L76)
``` go
type PreRootObjects struct {
    *Objects
}
```
PreRootObjects struct of *Objects for pre-configured hosts content










### <a name="PreRootObjects.Find">func</a> (\*PreRootObjects) [Find](/src/target/content.go?s=2673:2716#L151)
``` go
func (p *PreRootObjects) Find(s string) int
```
Find returns the int position of an Objects' element




### <a name="PreRootObjects.GetList">func</a> (\*PreRootObjects) [GetList](/src/target/content.go?s=4932:4975#L262)
``` go
func (p *PreRootObjects) GetList() *Objects
```
GetList implements the Contenter interface for PreRootObjects




### <a name="PreRootObjects.Len">func</a> (\*PreRootObjects) [Len](/src/target/content.go?s=6546:6580#L330)
``` go
func (p *PreRootObjects) Len() int
```
Len returns how many sources there are




### <a name="PreRootObjects.SetURL">func</a> (\*PreRootObjects) [SetURL](/src/target/content.go?s=9106:9155#L444)
``` go
func (p *PreRootObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="PreRootObjects.String">func</a> (\*PreRootObjects) [String](/src/target/content.go?s=9998:10038#L476)
``` go
func (p *PreRootObjects) String() string
```



## <a name="URLDomnObjects">type</a> [URLDomnObjects](/src/target/content.go?s=1360:1400#L81)
``` go
type URLDomnObjects struct {
    *Objects
}
```
URLDomnObjects struct of *Objects for domain URLs










### <a name="URLDomnObjects.Find">func</a> (\*URLDomnObjects) [Find](/src/target/content.go?s=3047:3090#L171)
``` go
func (u *URLDomnObjects) Find(s string) int
```
Find returns the int position of an Objects' element




### <a name="URLDomnObjects.GetList">func</a> (\*URLDomnObjects) [GetList](/src/target/content.go?s=5180:5223#L273)
``` go
func (u *URLDomnObjects) GetList() *Objects
```
GetList implements the Contenter interface for URLHostObjects




### <a name="URLDomnObjects.Len">func</a> (\*URLDomnObjects) [Len](/src/target/content.go?s=6646:6680#L333)
``` go
func (u *URLDomnObjects) Len() int
```
Len returns how many sources there are




### <a name="URLDomnObjects.SetURL">func</a> (\*URLDomnObjects) [SetURL](/src/target/content.go?s=9276:9325#L453)
``` go
func (u *URLDomnObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="URLDomnObjects.String">func</a> (\*URLDomnObjects) [String](/src/target/content.go?s=10069:10109#L477)
``` go
func (u *URLDomnObjects) String() string
```



## <a name="URLHostObjects">type</a> [URLHostObjects](/src/target/content.go?s=1453:1493#L86)
``` go
type URLHostObjects struct {
    *Objects
}
```
URLHostObjects struct of *Objects for host URLs










### <a name="URLHostObjects.Find">func</a> (\*URLHostObjects) [Find](/src/target/content.go?s=2860:2903#L161)
``` go
func (u *URLHostObjects) Find(s string) int
```
Find returns the int position of an Objects' element




### <a name="URLHostObjects.GetList">func</a> (\*URLHostObjects) [GetList](/src/target/content.go?s=5574:5617#L292)
``` go
func (u *URLHostObjects) GetList() *Objects
```
GetList implements the Contenter interface for URLHostObjects




### <a name="URLHostObjects.Len">func</a> (\*URLHostObjects) [Len](/src/target/content.go?s=6746:6780#L336)
``` go
func (u *URLHostObjects) Len() int
```
Len returns how many sources there are




### <a name="URLHostObjects.SetURL">func</a> (\*URLHostObjects) [SetURL](/src/target/content.go?s=9446:9495#L462)
``` go
func (u *URLHostObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="URLHostObjects.String">func</a> (\*URLHostObjects) [String](/src/target/content.go?s=10140:10180#L478)
``` go
func (u *URLHostObjects) String() string
```



## <a name="Wildcard">type</a> [Wildcard](/src/target/opts.go?s=1458:1556#L48)
``` go
type Wildcard struct {
    Node string `json:"Node,omitempty"`
    Name string `json:"Name,omitempty"`
}
```
Wildcard struct sets globbing wildcards for filename searches














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
