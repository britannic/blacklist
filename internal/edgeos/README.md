

# edgeos
`import "github.com/britannic/blacklist/internal/edgeos"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
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
  * [func (c *Config) InSession() bool](#Config.InSession)
  * [func (c *Config) LTypes() []string](#Config.LTypes)
  * [func (c *Config) NewContent(iface IFace) (Contenter, error)](#Config.NewContent)
  * [func (c *Config) Nodes() (nodes []string)](#Config.Nodes)
  * [func (c *Config) ProcessContent(cts ...Contenter) error](#Config.ProcessContent)
  * [func (c *Config) ReadCfg(r ConfLoader) error](#Config.ReadCfg)
  * [func (c *Config) ReloadDNS() ([]byte, error)](#Config.ReloadDNS)
  * [func (c *Config) SetOpt(opts ...Option) Option](#Config.SetOpt)
  * [func (c *Config) String() (s string)](#Config.String)
* [type Contenter](#Contenter)
* [type ExcDomnObjects](#ExcDomnObjects)
  * [func (e *ExcDomnObjects) Find(elem string) int](#ExcDomnObjects.Find)
  * [func (e *ExcDomnObjects) GetList() *Objects](#ExcDomnObjects.GetList)
  * [func (e *ExcDomnObjects) Len() int](#ExcDomnObjects.Len)
  * [func (e *ExcDomnObjects) SetURL(name, url string)](#ExcDomnObjects.SetURL)
  * [func (e *ExcDomnObjects) String() string](#ExcDomnObjects.String)
* [type ExcHostObjects](#ExcHostObjects)
  * [func (e *ExcHostObjects) Find(elem string) int](#ExcHostObjects.Find)
  * [func (e *ExcHostObjects) GetList() *Objects](#ExcHostObjects.GetList)
  * [func (e *ExcHostObjects) Len() int](#ExcHostObjects.Len)
  * [func (e *ExcHostObjects) SetURL(name, url string)](#ExcHostObjects.SetURL)
  * [func (e *ExcHostObjects) String() string](#ExcHostObjects.String)
* [type ExcRootObjects](#ExcRootObjects)
  * [func (e *ExcRootObjects) Find(elem string) int](#ExcRootObjects.Find)
  * [func (e *ExcRootObjects) GetList() *Objects](#ExcRootObjects.GetList)
  * [func (e *ExcRootObjects) Len() int](#ExcRootObjects.Len)
  * [func (e *ExcRootObjects) SetURL(name, url string)](#ExcRootObjects.SetURL)
  * [func (e *ExcRootObjects) String() string](#ExcRootObjects.String)
* [type FIODataObjects](#FIODataObjects)
  * [func (f *FIODataObjects) Find(elem string) int](#FIODataObjects.Find)
  * [func (f *FIODataObjects) GetList() *Objects](#FIODataObjects.GetList)
  * [func (f *FIODataObjects) Len() int](#FIODataObjects.Len)
  * [func (f *FIODataObjects) SetURL(name, url string)](#FIODataObjects.SetURL)
  * [func (f *FIODataObjects) String() string](#FIODataObjects.String)
* [type IFace](#IFace)
  * [func (i IFace) String() (s string)](#IFace.String)
* [type Msg](#Msg)
  * [func NewMsg(s string) *Msg](#NewMsg)
  * [func (m *Msg) GetTotal()](#Msg.GetTotal)
  * [func (m *Msg) IncDupe()](#Msg.IncDupe)
  * [func (m *Msg) IncNew()](#Msg.IncNew)
  * [func (m *Msg) IncUniq()](#Msg.IncUniq)
  * [func (m *Msg) String() string](#Msg.String)
* [type Objects](#Objects)
  * [func (o *Objects) Files() *CFile](#Objects.Files)
  * [func (o *Objects) Filter(ltype string) *Objects](#Objects.Filter)
  * [func (o *Objects) Find(elem string) int](#Objects.Find)
  * [func (o *Objects) Len() int](#Objects.Len)
  * [func (o *Objects) Less(i, j int) bool](#Objects.Less)
  * [func (o *Objects) Names() (s sort.StringSlice)](#Objects.Names)
  * [func (o *Objects) String() string](#Objects.String)
  * [func (o *Objects) Swap(i, j int)](#Objects.Swap)
* [type Option](#Option)
  * [func API(s string) Option](#API)
  * [func Arch(arch string) Option](#Arch)
  * [func Bash(cmd string) Option](#Bash)
  * [func Cores(i int) Option](#Cores)
  * [func DNSsvc(d string) Option](#DNSsvc)
  * [func Dbug(b bool) Option](#Dbug)
  * [func Dir(d string) Option](#Dir)
  * [func Ext(e string) Option](#Ext)
  * [func File(f string) Option](#File)
  * [func FileNameFmt(f string) Option](#FileNameFmt)
  * [func InCLI(in string) Option](#InCLI)
  * [func LTypes(s []string) Option](#LTypes)
  * [func Level(s string) Option](#Level)
  * [func Logger(l *logging.Logger) Option](#Logger)
  * [func Method(method string) Option](#Method)
  * [func Nodes(nodes []string) Option](#Nodes)
  * [func Poll(t int) Option](#Poll)
  * [func Prefix(l string) Option](#Prefix)
  * [func Test(b bool) Option](#Test)
  * [func Timeout(t time.Duration) Option](#Timeout)
  * [func Verb(b bool) Option](#Verb)
  * [func WCard(w Wildcard) Option](#WCard)
  * [func Writer(w io.Writer) Option](#Writer)
* [type Parms](#Parms)
  * [func (p *Parms) String() string](#Parms.String)
* [type PreDomnObjects](#PreDomnObjects)
  * [func (p *PreDomnObjects) Find(elem string) int](#PreDomnObjects.Find)
  * [func (p *PreDomnObjects) GetList() *Objects](#PreDomnObjects.GetList)
  * [func (p *PreDomnObjects) Len() int](#PreDomnObjects.Len)
  * [func (p *PreDomnObjects) SetURL(name, url string)](#PreDomnObjects.SetURL)
  * [func (p *PreDomnObjects) String() string](#PreDomnObjects.String)
* [type PreHostObjects](#PreHostObjects)
  * [func (p *PreHostObjects) Find(elem string) int](#PreHostObjects.Find)
  * [func (p *PreHostObjects) GetList() *Objects](#PreHostObjects.GetList)
  * [func (p *PreHostObjects) Len() int](#PreHostObjects.Len)
  * [func (p *PreHostObjects) SetURL(name, url string)](#PreHostObjects.SetURL)
  * [func (p *PreHostObjects) String() string](#PreHostObjects.String)
* [type Rec](#Rec)
* [type URLDomnObjects](#URLDomnObjects)
  * [func (u *URLDomnObjects) Find(elem string) int](#URLDomnObjects.Find)
  * [func (u *URLDomnObjects) GetList() *Objects](#URLDomnObjects.GetList)
  * [func (u *URLDomnObjects) Len() int](#URLDomnObjects.Len)
  * [func (u *URLDomnObjects) SetURL(name, url string)](#URLDomnObjects.SetURL)
  * [func (u *URLDomnObjects) String() string](#URLDomnObjects.String)
* [type URLHostObjects](#URLHostObjects)
  * [func (u *URLHostObjects) Find(elem string) int](#URLHostObjects.Find)
  * [func (u *URLHostObjects) GetList() *Objects](#URLHostObjects.GetList)
  * [func (u *URLHostObjects) Len() int](#URLHostObjects.Len)
  * [func (u *URLHostObjects) SetURL(name, url string)](#URLHostObjects.SetURL)
  * [func (u *URLHostObjects) String() string](#URLHostObjects.String)
* [type Wildcard](#Wildcard)


#### <a name="pkg-files">Package files</a>
[config.go](/src/github.com/britannic/blacklist/internal/edgeos/config.go) [content.go](/src/github.com/britannic/blacklist/internal/edgeos/content.go) [data.go](/src/github.com/britannic/blacklist/internal/edgeos/data.go) [http.go](/src/github.com/britannic/blacklist/internal/edgeos/http.go) [io.go](/src/github.com/britannic/blacklist/internal/edgeos/io.go) [json.go](/src/github.com/britannic/blacklist/internal/edgeos/json.go) [list.go](/src/github.com/britannic/blacklist/internal/edgeos/list.go) [msg.go](/src/github.com/britannic/blacklist/internal/edgeos/msg.go) [ntype_string.go](/src/github.com/britannic/blacklist/internal/edgeos/ntype_string.go) [object.go](/src/github.com/britannic/blacklist/internal/edgeos/object.go) [opts.go](/src/github.com/britannic/blacklist/internal/edgeos/opts.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (

    // ExcDomns labels domain exclusions
    ExcDomns = "domn-excludes"
    // ExcHosts labels host exclusions
    ExcHosts = "host-excludes"
    // ExcRoots labels global domain exclusions
    ExcRoots = "root-excludes"
    // PreDomns designates string label for preconfigured blacklisted domains
    PreDomns = preNoun + "-domain"
    // PreHosts designates string label for preconfigured blacklisted hosts
    PreHosts = preNoun + "-host"
    // False is a string constant
    False = "false"
    // True is a string constant
    True = "true"
)
```



## <a name="Iter">func</a> [Iter](/src/target/data.go?s=2391:2418#L99)
``` go
func Iter(i int) []struct{}
```
Iter iterates over ints - use it in for loops



## <a name="NewWriter">func</a> [NewWriter](/src/target/data.go?s=2486:2512#L104)
``` go
func NewWriter() io.Writer
```
NewWriter returns an io.Writer




## <a name="CFGcli">type</a> [CFGcli](/src/target/io.go?s=139:182#L4)
``` go
type CFGcli struct {
    *Config
    Cfg string
}
```
CFGcli loads configurations using the EdgeOS CFGcli










## <a name="CFGstatic">type</a> [CFGstatic](/src/target/io.go?s=237:283#L10)
``` go
type CFGstatic struct {
    *Config
    Cfg string
}
```
CFGstatic loads static configurations for testing










## <a name="CFile">type</a> [CFile](/src/target/config.go?s=487:545#L17)
``` go
type CFile struct {
    *Parms
    // contains filtered or unexported fields
}
```
CFile holds an array of file names










### <a name="CFile.Remove">func</a> (\*CFile) [Remove](/src/target/config.go?s=8027:8057#L365)
``` go
func (c *CFile) Remove() error
```
Remove deletes a CFile array of file names




### <a name="CFile.String">func</a> (\*CFile) [String](/src/target/config.go?s=9444:9475#L424)
``` go
func (c *CFile) String() string
```
String implements string method




### <a name="CFile.Strings">func</a> (\*CFile) [Strings](/src/target/config.go?s=9563:9597#L429)
``` go
func (c *CFile) Strings() []string
```
Strings returns a sorted array of strings.




## <a name="ConfLoader">type</a> [ConfLoader](/src/target/config.go?s=400:447#L12)
``` go
type ConfLoader interface {
    // contains filtered or unexported methods
}
```
ConfLoader interface defines configuration load method










## <a name="Config">type</a> [Config](/src/target/config.go?s=593:629#L24)
``` go
type Config struct {
    *Parms
    // contains filtered or unexported fields
}
```
Config is a struct of configuration fields







### <a name="NewConfig">func</a> [NewConfig](/src/target/opts.go?s=4850:4888#L207)
``` go
func NewConfig(opts ...Option) *Config
```
NewConfig returns a new *Config initialized with the parameter options passed to it





### <a name="Config.Get">func</a> (\*Config) [Get](/src/target/config.go?s=4012:4054#L185)
``` go
func (c *Config) Get(node string) *Objects
```
Get returns an *Object for a given node




### <a name="Config.GetAll">func</a> (\*Config) [GetAll](/src/target/config.go?s=4279:4329#L200)
``` go
func (c *Config) GetAll(ltypes ...string) *Objects
```
GetAll returns an array of Objects




### <a name="Config.InSession">func</a> (\*Config) [InSession](/src/target/config.go?s=5067:5100#L239)
``` go
func (c *Config) InSession() bool
```
InSession returns true if VyOS/EdgeOS configuration is in session




### <a name="Config.LTypes">func</a> (\*Config) [LTypes](/src/target/config.go?s=9689:9723#L435)
``` go
func (c *Config) LTypes() []string
```
LTypes returns an array of configured nodes




### <a name="Config.NewContent">func</a> (\*Config) [NewContent](/src/target/config.go?s=2628:2687#L119)
``` go
func (c *Config) NewContent(iface IFace) (Contenter, error)
```
NewContent returns an interface of the requested IFace type




### <a name="Config.Nodes">func</a> (\*Config) [Nodes](/src/target/config.go?s=5468:5509#L252)
``` go
func (c *Config) Nodes() (nodes []string)
```
Nodes returns an array of configured nodes




### <a name="Config.ProcessContent">func</a> (\*Config) [ProcessContent](/src/target/content.go?s=7705:7760#L377)
``` go
func (c *Config) ProcessContent(cts ...Contenter) error
```
ProcessContent processes the Contents array




### <a name="Config.ReadCfg">func</a> (\*Config) [ReadCfg](/src/target/config.go?s=5675:5719#L262)
``` go
func (c *Config) ReadCfg(r ConfLoader) error
```
ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure




### <a name="Config.ReloadDNS">func</a> (\*Config) [ReloadDNS](/src/target/config.go?s=7831:7875#L357)
``` go
func (c *Config) ReloadDNS() ([]byte, error)
```
ReloadDNS reloads the dnsmasq configuration




### <a name="Config.SetOpt">func</a> (\*Config) [SetOpt](/src/target/opts.go?s=1853:1899#L55)
``` go
func (c *Config) SetOpt(opts ...Option) Option
```
SetOpt sets the specified options passed as Parms and returns an option to restore the last set of arg's previous values




### <a name="Config.String">func</a> (\*Config) [String](/src/target/config.go?s=8550:8586#L388)
``` go
func (c *Config) String() (s string)
```
String returns pretty print for the Blacklist struct




## <a name="Contenter">type</a> [Contenter](/src/target/content.go?s=436:554#L27)
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










## <a name="ExcDomnObjects">type</a> [ExcDomnObjects](/src/target/content.go?s=704:744#L41)
``` go
type ExcDomnObjects struct {
    *Objects
}
```
ExcDomnObjects implements GetList for domain exclusions










### <a name="ExcDomnObjects.Find">func</a> (\*ExcDomnObjects) [Find](/src/target/content.go?s=1411:1457#L76)
``` go
func (e *ExcDomnObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="ExcDomnObjects.GetList">func</a> (\*ExcDomnObjects) [GetList](/src/target/content.go?s=2900:2943#L156)
``` go
func (e *ExcDomnObjects) GetList() *Objects
```
GetList implements the Contenter interface for ExcDomnObjects




### <a name="ExcDomnObjects.Len">func</a> (\*ExcDomnObjects) [Len](/src/target/content.go?s=5578:5612#L290)
``` go
func (e *ExcDomnObjects) Len() int
```
Len returns how many objects there are




### <a name="ExcDomnObjects.SetURL">func</a> (\*ExcDomnObjects) [SetURL](/src/target/content.go?s=8587:8636#L425)
``` go
func (e *ExcDomnObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="ExcDomnObjects.String">func</a> (\*ExcDomnObjects) [String](/src/target/content.go?s=9887:9927#L496)
``` go
func (e *ExcDomnObjects) String() string
```



## <a name="ExcHostObjects">type</a> [ExcHostObjects](/src/target/content.go?s=803:843#L46)
``` go
type ExcHostObjects struct {
    *Objects
}
```
ExcHostObjects implements GetList for host exclusions










### <a name="ExcHostObjects.Find">func</a> (\*ExcHostObjects) [Find](/src/target/content.go?s=1596:1642#L86)
``` go
func (e *ExcHostObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="ExcHostObjects.GetList">func</a> (\*ExcHostObjects) [GetList](/src/target/content.go?s=3178:3221#L170)
``` go
func (e *ExcHostObjects) GetList() *Objects
```
GetList implements the Contenter interface for ExcHostObjects




### <a name="ExcHostObjects.Len">func</a> (\*ExcHostObjects) [Len](/src/target/content.go?s=5684:5718#L293)
``` go
func (e *ExcHostObjects) Len() int
```
Len returns how many objects there are




### <a name="ExcHostObjects.SetURL">func</a> (\*ExcHostObjects) [SetURL](/src/target/content.go?s=8755:8804#L434)
``` go
func (e *ExcHostObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="ExcHostObjects.String">func</a> (\*ExcHostObjects) [String](/src/target/content.go?s=9958:9998#L497)
``` go
func (e *ExcHostObjects) String() string
```



## <a name="ExcRootObjects">type</a> [ExcRootObjects](/src/target/content.go?s=911:951#L51)
``` go
type ExcRootObjects struct {
    *Objects
}
```
ExcRootObjects implements GetList for global domain exclusions










### <a name="ExcRootObjects.Find">func</a> (\*ExcRootObjects) [Find](/src/target/content.go?s=1781:1827#L96)
``` go
func (e *ExcRootObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="ExcRootObjects.GetList">func</a> (\*ExcRootObjects) [GetList](/src/target/content.go?s=3456:3499#L184)
``` go
func (e *ExcRootObjects) GetList() *Objects
```
GetList implements the Contenter interface for ExcRootObjects




### <a name="ExcRootObjects.Len">func</a> (\*ExcRootObjects) [Len](/src/target/content.go?s=5790:5824#L296)
``` go
func (e *ExcRootObjects) Len() int
```
Len returns how many objects there are




### <a name="ExcRootObjects.SetURL">func</a> (\*ExcRootObjects) [SetURL](/src/target/content.go?s=8923:8972#L443)
``` go
func (e *ExcRootObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="ExcRootObjects.String">func</a> (\*ExcRootObjects) [String](/src/target/content.go?s=10029:10069#L498)
``` go
func (e *ExcRootObjects) String() string
```



## <a name="FIODataObjects">type</a> [FIODataObjects](/src/target/content.go?s=603:643#L36)
``` go
type FIODataObjects struct {
    *Objects
}
```
FIODataObjects implements GetList for files










### <a name="FIODataObjects.Find">func</a> (\*FIODataObjects) [Find](/src/target/content.go?s=1966:2012#L106)
``` go
func (f *FIODataObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="FIODataObjects.GetList">func</a> (\*FIODataObjects) [GetList](/src/target/content.go?s=3734:3777#L198)
``` go
func (f *FIODataObjects) GetList() *Objects
```
GetList implements the Contenter interface for FIODataObjects




### <a name="FIODataObjects.Len">func</a> (\*FIODataObjects) [Len](/src/target/content.go?s=5896:5930#L299)
``` go
func (f *FIODataObjects) Len() int
```
Len returns how many objects there are




### <a name="FIODataObjects.SetURL">func</a> (\*FIODataObjects) [SetURL](/src/target/content.go?s=9091:9140#L452)
``` go
func (f *FIODataObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="FIODataObjects.String">func</a> (\*FIODataObjects) [String](/src/target/content.go?s=10100:10140#L499)
``` go
func (f *FIODataObjects) String() string
```



## <a name="IFace">type</a> [IFace](/src/target/content.go?s=177:191#L6)
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
    URLdObj
    URLhObj
)
```
IFace types for labeling interface types










### <a name="IFace.String">func</a> (IFace) [String](/src/target/content.go?s=10456:10490#L505)
``` go
func (i IFace) String() (s string)
```



## <a name="Msg">type</a> [Msg](/src/target/msg.go?s=287:325#L8)
``` go
type Msg struct {
    Name string
    *Rec
}
```
Msg is a struct for recording stats from ProcessContent







### <a name="NewMsg">func</a> [NewMsg](/src/target/msg.go?s=820:846#L49)
``` go
func NewMsg(s string) *Msg
```
NewMsg initializes a new Msg struct





### <a name="Msg.GetTotal">func</a> (\*Msg) [GetTotal](/src/target/msg.go?s=435:459#L21)
``` go
func (m *Msg) GetTotal()
```
GetTotal returns total records




### <a name="Msg.IncDupe">func</a> (\*Msg) [IncDupe](/src/target/msg.go?s=537:560#L28)
``` go
func (m *Msg) IncDupe()
```
IncDupe increments Dupe by 1




### <a name="Msg.IncNew">func</a> (\*Msg) [IncNew](/src/target/msg.go?s=629:651#L35)
``` go
func (m *Msg) IncNew()
```
IncNew increments New by 1




### <a name="Msg.IncUniq">func</a> (\*Msg) [IncUniq](/src/target/msg.go?s=720:743#L42)
``` go
func (m *Msg) IncUniq()
```
IncUniq increments Uniq by 1




### <a name="Msg.String">func</a> (\*Msg) [String](/src/target/msg.go?s=927:956#L58)
``` go
func (m *Msg) String() string
```



## <a name="Objects">type</a> [Objects](/src/target/object.go?s=401:445#L20)
``` go
type Objects struct {
    *Parms
    // contains filtered or unexported fields
}
```
Objects is a struct of []*Object










### <a name="Objects.Files">func</a> (\*Objects) [Files](/src/target/object.go?s=915:947#L42)
``` go
func (o *Objects) Files() *CFile
```
Files returns a list of dnsmasq conf files from all srcs




### <a name="Objects.Filter">func</a> (\*Objects) [Filter](/src/target/object.go?s=1274:1321#L54)
``` go
func (o *Objects) Filter(ltype string) *Objects
```
Filter returns a subset of Objects; ltypes with "-" prepended remove ltype




### <a name="Objects.Find">func</a> (\*Objects) [Find](/src/target/object.go?s=2056:2095#L93)
``` go
func (o *Objects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="Objects.Len">func</a> (\*Objects) [Len](/src/target/object.go?s=3309:3336#L146)
``` go
func (o *Objects) Len() int
```
Implement Sort Interface for Objects




### <a name="Objects.Less">func</a> (\*Objects) [Less](/src/target/object.go?s=3367:3404#L147)
``` go
func (o *Objects) Less(i, j int) bool
```



### <a name="Objects.Names">func</a> (\*Objects) [Names](/src/target/object.go?s=2403:2449#L109)
``` go
func (o *Objects) Names() (s sort.StringSlice)
```
Names returns a sorted slice of Objects names




### <a name="Objects.String">func</a> (\*Objects) [String](/src/target/object.go?s=3206:3239#L141)
``` go
func (o *Objects) String() string
```
Stringer for Objects




### <a name="Objects.Swap">func</a> (\*Objects) [Swap](/src/target/object.go?s=3442:3474#L148)
``` go
func (o *Objects) Swap(i, j int)
```



## <a name="Option">type</a> [Option](/src/target/opts.go?s=1548:1582#L39)
``` go
type Option func(c *Config) Option
```
Option is a recursive function







### <a name="API">func</a> [API](/src/target/opts.go?s=2385:2410#L80)
``` go
func API(s string) Option
```
API sets the EdgeOS CLI API command


### <a name="Arch">func</a> [Arch](/src/target/opts.go?s=2214:2243#L71)
``` go
func Arch(arch string) Option
```
Arch sets target CPU architecture


### <a name="Bash">func</a> [Bash](/src/target/opts.go?s=2540:2568#L89)
``` go
func Bash(cmd string) Option
```
Bash sets the shell processor


### <a name="Cores">func</a> [Cores](/src/target/opts.go?s=2698:2722#L98)
``` go
func Cores(i int) Option
```
Cores sets max CPU cores


### <a name="DNSsvc">func</a> [DNSsvc](/src/target/opts.go?s=3203:3231#L126)
``` go
func DNSsvc(d string) Option
```
DNSsvc sets dnsmasq restart command


### <a name="Dbug">func</a> [Dbug](/src/target/opts.go?s=2887:2911#L108)
``` go
func Dbug(b bool) Option
```
Dbug toggles debug level on or off


### <a name="Dir">func</a> [Dir](/src/target/opts.go?s=3042:3067#L117)
``` go
func Dir(d string) Option
```
Dir sets directory location


### <a name="Ext">func</a> [Ext](/src/target/opts.go?s=3380:3405#L135)
``` go
func Ext(e string) Option
```
Ext sets the blacklist file n extension


### <a name="File">func</a> [File](/src/target/opts.go?s=3545:3571#L144)
``` go
func File(f string) Option
```
File sets the EdgeOS configuration file


### <a name="FileNameFmt">func</a> [FileNameFmt](/src/target/opts.go?s=3733:3766#L153)
``` go
func FileNameFmt(f string) Option
```
FileNameFmt sets the EdgeOS configuration file name format


### <a name="InCLI">func</a> [InCLI](/src/target/opts.go?s=3915:3943#L162)
``` go
func InCLI(in string) Option
```
InCLI sets the CLI inSession command


### <a name="LTypes">func</a> [LTypes](/src/target/opts.go?s=4452:4482#L189)
``` go
func LTypes(s []string) Option
```
LTypes sets an array of legal types used by Source


### <a name="Level">func</a> [Level](/src/target/opts.go?s=4086:4113#L171)
``` go
func Level(s string) Option
```
Level sets the EdgeOS API CLI level


### <a name="Logger">func</a> [Logger](/src/target/opts.go?s=4255:4292#L180)
``` go
func Logger(l *logging.Logger) Option
```
Logger sets a pointer to the logger


### <a name="Method">func</a> [Method](/src/target/opts.go?s=4619:4652#L198)
``` go
func Method(method string) Option
```
Method sets the HTTP method


### <a name="Nodes">func</a> [Nodes](/src/target/opts.go?s=5157:5190#L222)
``` go
func Nodes(nodes []string) Option
```
Nodes sets the node ns array


### <a name="Poll">func</a> [Poll](/src/target/opts.go?s=5354:5377#L231)
``` go
func Poll(t int) Option
```
Poll sets the polling interval in seconds


### <a name="Prefix">func</a> [Prefix](/src/target/opts.go?s=5538:5566#L240)
``` go
func Prefix(l string) Option
```
Prefix sets the dnsmasq configuration address line prefix


### <a name="Test">func</a> [Test](/src/target/opts.go?s=5855:5879#L255)
``` go
func Test(b bool) Option
```
Test toggles testing mode on or off


### <a name="Timeout">func</a> [Timeout](/src/target/opts.go?s=6048:6084#L264)
``` go
func Timeout(t time.Duration) Option
```
Timeout sets how long before an unresponsive goroutine is aborted


### <a name="Verb">func</a> [Verb](/src/target/opts.go?s=6231:6255#L273)
``` go
func Verb(b bool) Option
```
Verb sets the verbosity level to v


### <a name="WCard">func</a> [WCard](/src/target/opts.go?s=6399:6428#L282)
``` go
func WCard(w Wildcard) Option
```
WCard sets file globbing wildcard values


### <a name="Writer">func</a> [Writer](/src/target/opts.go?s=6603:6634#L291)
``` go
func Writer(w io.Writer) Option
```
Writer provides an address for anything that can use io.Writer





## <a name="Parms">type</a> [Parms](/src/target/opts.go?s=148:1357#L4)
``` go
type Parms struct {
    *logging.Logger
    API      string        `json:"API, omitempty"`
    Arch     string        `json:"Arch, omitempty"`
    Bash     string        `json:"Bash, omitempty"`
    Cores    int           `json:"Cores, omitempty"`
    Dbug     bool          `json:"Dbug, omitempty"`
    Dex      list          `json:"Dex, omitempty"`
    Dir      string        `json:"Dir, omitempty"`
    DNSsvc   string        `json:"dnsmasq service, omitempty"`
    Exc      list          `json:"Exc, omitempty"`
    Ext      string        `json:"dnsmasq fileExt., omitempty"`
    File     string        `json:"File, omitempty"`
    FnFmt    string        `json:"File name fmt, omitempty"`
    InCLI    string        `json:"-"`
    Level    string        `json:"CLI Path, omitempty"`
    Ltypes   []string      `json:"Leaf nodes, omitempty"`
    Method   string        `json:"HTTP method, omitempty"`
    Nodes    []string      `json:"Nodes, omitempty"`
    Pfx      string        `json:"Prefix, omitempty"`
    Poll     int           `json:"Poll, omitempty"`
    Test     bool          `json:"Test, omitempty"`
    Timeout  time.Duration `json:"Timeout, omitempty"`
    Verb     bool          `json:"Verbosity, omitempty"`
    Wildcard `json:"Wildcard, omitempty"`
    // contains filtered or unexported fields
}
```
Parms is struct of parameters










### <a name="Parms.String">func</a> (\*Parms) [String](/src/target/opts.go?s=5716:5747#L249)
``` go
func (p *Parms) String() string
```
String method to implement fmt.Print interface




## <a name="PreDomnObjects">type</a> [PreDomnObjects](/src/target/content.go?s=1025:1065#L56)
``` go
type PreDomnObjects struct {
    *Objects
}
```
PreDomnObjects implements GetList for pre-configured domains content










### <a name="PreDomnObjects.Find">func</a> (\*PreDomnObjects) [Find](/src/target/content.go?s=2151:2197#L116)
``` go
func (p *PreDomnObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="PreDomnObjects.GetList">func</a> (\*PreDomnObjects) [GetList](/src/target/content.go?s=4199:4242#L222)
``` go
func (p *PreDomnObjects) GetList() *Objects
```
GetList implements the Contenter interface for PreDomnObjects




### <a name="PreDomnObjects.Len">func</a> (\*PreDomnObjects) [Len](/src/target/content.go?s=6002:6036#L302)
``` go
func (p *PreDomnObjects) Len() int
```
Len returns how many objects there are




### <a name="PreDomnObjects.SetURL">func</a> (\*PreDomnObjects) [SetURL](/src/target/content.go?s=9259:9308#L461)
``` go
func (p *PreDomnObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="PreDomnObjects.String">func</a> (\*PreDomnObjects) [String](/src/target/content.go?s=10171:10211#L500)
``` go
func (p *PreDomnObjects) String() string
```



## <a name="PreHostObjects">type</a> [PreHostObjects](/src/target/content.go?s=1137:1177#L61)
``` go
type PreHostObjects struct {
    *Objects
}
```
PreHostObjects implements GetList for pre-configured hosts content










### <a name="PreHostObjects.Find">func</a> (\*PreHostObjects) [Find](/src/target/content.go?s=2336:2382#L126)
``` go
func (p *PreHostObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="PreHostObjects.GetList">func</a> (\*PreHostObjects) [GetList](/src/target/content.go?s=4457:4500#L233)
``` go
func (p *PreHostObjects) GetList() *Objects
```
GetList implements the Contenter interface for PreHostObjects




### <a name="PreHostObjects.Len">func</a> (\*PreHostObjects) [Len](/src/target/content.go?s=6108:6142#L305)
``` go
func (p *PreHostObjects) Len() int
```
Len returns how many objects there are




### <a name="PreHostObjects.SetURL">func</a> (\*PreHostObjects) [SetURL](/src/target/content.go?s=9427:9476#L470)
``` go
func (p *PreHostObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="PreHostObjects.String">func</a> (\*PreHostObjects) [String](/src/target/content.go?s=10242:10282#L501)
``` go
func (p *PreHostObjects) String() string
```



## <a name="Rec">type</a> [Rec](/src/target/msg.go?s=91:226#L1)
``` go
type Rec struct {
    *sync.RWMutex
    Dupes int `json:"dupes"`
    New   int `json:"new"`
    Total int `json:"total"`
    Uniq  int `json:"uniq"`
}
```
Rec holds stats on the current job










## <a name="URLDomnObjects">type</a> [URLDomnObjects](/src/target/content.go?s=1313:1353#L71)
``` go
type URLDomnObjects struct {
    *Objects
}
```
URLDomnObjects implements GetList for URLs










### <a name="URLDomnObjects.Find">func</a> (\*URLDomnObjects) [Find](/src/target/content.go?s=2706:2752#L146)
``` go
func (u *URLDomnObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="URLDomnObjects.GetList">func</a> (\*URLDomnObjects) [GetList](/src/target/content.go?s=4715:4758#L244)
``` go
func (u *URLDomnObjects) GetList() *Objects
```
GetList implements the Contenter interface for URLHostObjects




### <a name="URLDomnObjects.Len">func</a> (\*URLDomnObjects) [Len](/src/target/content.go?s=6214:6248#L308)
``` go
func (u *URLDomnObjects) Len() int
```
Len returns how many objects there are




### <a name="URLDomnObjects.SetURL">func</a> (\*URLDomnObjects) [SetURL](/src/target/content.go?s=9595:9644#L479)
``` go
func (u *URLDomnObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="URLDomnObjects.String">func</a> (\*URLDomnObjects) [String](/src/target/content.go?s=10313:10353#L502)
``` go
func (u *URLDomnObjects) String() string
```



## <a name="URLHostObjects">type</a> [URLHostObjects](/src/target/content.go?s=1225:1265#L66)
``` go
type URLHostObjects struct {
    *Objects
}
```
URLHostObjects implements GetList for URLs










### <a name="URLHostObjects.Find">func</a> (\*URLHostObjects) [Find](/src/target/content.go?s=2521:2567#L136)
``` go
func (u *URLHostObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="URLHostObjects.GetList">func</a> (\*URLHostObjects) [GetList](/src/target/content.go?s=5158:5201#L267)
``` go
func (u *URLHostObjects) GetList() *Objects
```
GetList implements the Contenter interface for URLHostObjects




### <a name="URLHostObjects.Len">func</a> (\*URLHostObjects) [Len](/src/target/content.go?s=6320:6354#L311)
``` go
func (u *URLHostObjects) Len() int
```
Len returns how many objects there are




### <a name="URLHostObjects.SetURL">func</a> (\*URLHostObjects) [SetURL](/src/target/content.go?s=9763:9812#L488)
``` go
func (u *URLHostObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="URLHostObjects.String">func</a> (\*URLHostObjects) [String](/src/target/content.go?s=10384:10424#L503)
``` go
func (u *URLHostObjects) String() string
```



## <a name="Wildcard">type</a> [Wildcard](/src/target/opts.go?s=1424:1512#L33)
``` go
type Wildcard struct {
    Node string `json:"omitempty"`
    Name string `json:"omitempty"`
}
```
Wildcard struct sets globbing wildcards for filename searches














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
