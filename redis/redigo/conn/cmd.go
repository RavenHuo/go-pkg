/**
 * @Author raven
 * @Description
 * @Date 2022/7/27
 **/
package conn

import (
	"strings"
	"time"

	"github.com/RavenHuo/go-kit/internal"
)

type ICmder interface {
	Cmd() string
	Args() []interface{}
	Err() error
	SetErr(err error)
	ToString() string
}
type baseCmder struct {
	cmd  string
	args []interface{}
	err  error
	val  interface{}
}

func (c *baseCmder) Val() interface{} {
	return c.val
}

func (c *baseCmder) setVal(v interface{}) {
	c.val = v
}

func (c *baseCmder) Result() (interface{}, error) {
	return c.val, c.err
}

func (c *baseCmder) Err() error {
	return c.err
}
func (c *baseCmder) Args() []interface{} {
	return c.args
}
func (c *baseCmder) ToString() string {
	return cmdString(c, nil)
}
func (c *baseCmder) Cmd() string {
	return c.cmd
}
func (c *baseCmder) SetErr(err error) {
	c.err = err
}
func newBaseCmder(cmd string, arg ...interface{}) *baseCmder {
	return &baseCmder{
		cmd:  cmd,
		args: arg,
	}
}

func cmdString(c ICmder, val interface{}) string {
	b := make([]byte, 0, 64)

	b = internal.AppendArg(b, c.Cmd())
	b = append(b, ' ')

	for i, arg := range c.Args() {
		if i > 0 {
			b = append(b, ' ')
		}
		b = internal.AppendArg(b, arg)
	}

	if err := c.Err(); err != nil {
		b = append(b, ": "...)
		b = append(b, err.Error()...)
	} else if val != nil {
		b = append(b, ": "...)
		b = internal.AppendArg(b, val)
	}

	return internal.String(b)
}

///////////////////////////////// ZMember /////////////////////////
type ZMember struct {
	Score  float64
	Member interface{}
}
type ZMemberRangeBy struct {
	Min, Max      string
	Offset, Count int64
}

///////////////////////////////// ZMember /////////////////////////

///////////////////////////////// IntCmder ////////////////////////
type IntCmder struct {
	*baseCmder
	val int64
}

func (c *IntCmder) Val() int64 {
	return c.val
}

func (c *IntCmder) setVal(v int64) {
	c.val = v
}

func (c *IntCmder) Result() (int64, error) {
	return c.val, c.err
}

func (c *IntCmder) Uint64() (uint64, error) {
	return uint64(c.val), c.err
}

func (c *IntCmder) ToString() string {
	return cmdString(c.baseCmder, c.val)
}
func wrapperIntCmder(cmd *baseCmder) ICmder {
	intCmder := &IntCmder{
		baseCmder: cmd,
	}
	v, _ := cmd.Val().(int64)
	intCmder.setVal(v)
	return intCmder
}

///////////////////////////////// IntCmder /////////////////////////

///////////////////////////////// StringCmder //////////////////////
type StringCmder struct {
	*baseCmder
	val string
}

func (c *StringCmder) setVal(v string) {
	c.val = v
}

func (c *StringCmder) Val() string {
	return c.val
}

func (c *StringCmder) Result() (string, error) {
	return c.val, c.err
}

func (c *StringCmder) ToString() string {
	return cmdString(c.baseCmder, c.val)
}

func wrapperStringCmder(cmd *baseCmder) ICmder {
	s := &StringCmder{
		baseCmder: cmd,
	}
	v, _ := cmd.Val().([]byte)
	stringBuilder := strings.Builder{}
	stringBuilder.Write(v)
	s.setVal(stringBuilder.String())
	return s
}

///////////////////////////////// StringCmder ///////////////////////

///////////////////////////////// BoolCmder /////////////////////////
type BoolCmder struct {
	*baseCmder
	val bool
	response string
}

func (c *BoolCmder) setVal(v bool) {
	c.val = v
}

func (c *BoolCmder) Val() bool {
	return c.val
}

func (c *BoolCmder) Result() (bool, error) {
	return c.val, c.err
}

func (c *BoolCmder) ToString() string {
	return cmdString(c.baseCmder, c.val)
}
func (c *BoolCmder) Response() string {
	return c.response
}
func wrapperBoolCmder(cmd *baseCmder) ICmder {
	boolCmd := &BoolCmder{
		baseCmder: cmd,
	}
	v, _ := cmd.Val().(string)
	if v == oKResponse {
		boolCmd.setVal(true)
		boolCmd.response = v
	}
	return boolCmd
}

///////////////////////////////// BoolCmder /////////////////////////

///////////////////////////////// DurationCmder /////////////////////
type DurationCmder struct {
	*baseCmder
	val time.Duration
}

func (c *DurationCmder) setVal(v time.Duration) {
	c.val = v
}

func (c *DurationCmder) Val() time.Duration {
	return c.val
}

func (c *DurationCmder) Result() (time.Duration, error) {
	return c.val, c.err
}

func (c *DurationCmder) ToString() string {
	return cmdString(c.baseCmder, c.val)
}

func wrapperDurationCmder(cmd *baseCmder) ICmder {
	durationCmder := &DurationCmder{
		baseCmder: cmd,
	}
	seconds, _ := cmd.Val().(int64)
	v := time.Second * time.Duration(seconds)
	durationCmder.setVal(v)
	return durationCmder
}

///////////////////////////////// DurationCmder /////////////////////

///////////////////////////////// FloatCmder    /////////////////////
type FloatCmder struct {
	*baseCmder
	val float64
}

func (c *FloatCmder) setVal(v float64) {
	c.val = v
}

func (c *FloatCmder) Val() float64 {
	return c.val
}

func (c *FloatCmder) Result() (float64, error) {
	return c.val, c.err
}

func (c *FloatCmder) ToString() string {
	return cmdString(c.baseCmder, c.val)
}
func wrapperFloatCmder(cmd *baseCmder) ICmder {
	floatCmd := &FloatCmder{
		baseCmder: cmd,
	}
	v, _ := cmd.Val().(float64)
	floatCmd.setVal(v)
	return floatCmd
}

///////////////////////////////// FloatCmder    /////////////////////

///////////////////////////////// IntSliceCmder /////////////////////
type IntSliceCmder struct {
	*baseCmder
	val []int64
}

func (c *IntSliceCmder) setVal(v []int64) {
	c.val = v
}

func (c *IntSliceCmder) Val() []int64 {
	return c.val
}

func (c *IntSliceCmder) Result() ([]int64, error) {
	return c.val, c.err
}

func (c *IntSliceCmder) ToString() string {
	return cmdString(c.baseCmder, c.val)
}
func wrapperIntSliceCmder(cmd *baseCmder) ICmder {
	floatCmd := &IntSliceCmder{
		baseCmder: cmd,
	}
	v, _ := cmd.Val().([]int64)
	floatCmd.setVal(v)
	return floatCmd
}

///////////////////////////////// IntSliceCmder /////////////////////

///////////////////////////////// StringSliceCmder //////////////////
type StringSliceCmder struct {
	*baseCmder
	val []string
}

func (c *StringSliceCmder) setVal(v []string) {
	c.val = v
}

func (c *StringSliceCmder) Val() []string {
	return c.val
}

func (c *StringSliceCmder) Result() ([]string, error) {
	return c.val, c.err
}

func (c *StringSliceCmder) ToString() string {
	return cmdString(c.baseCmder, c.val)
}
func wrapperStringSliceCmder(cmd *baseCmder) ICmder {
	stringSliceCmd := &StringSliceCmder{
		baseCmder: cmd,
	}
	v, _ := cmd.Val().([]string)
	stringSliceCmd.setVal(v)
	return stringSliceCmd
}

///////////////////////////////// StringSliceCmder //////////////////

///////////////////////////////// InterfaceSliceCmder ///////////////
type InterfaceSliceCmder struct {
	*baseCmder
	val []interface{}
}

func (c *InterfaceSliceCmder) setVal(v []interface{}) {
	c.val = v
}

func (c *InterfaceSliceCmder) Val() []interface{} {
	return c.val
}

func (c *InterfaceSliceCmder) Result() ([]interface{}, error) {
	return c.val, c.err
}

func (c *InterfaceSliceCmder) ToString() string {
	return cmdString(c.baseCmder, c.val)
}
func wrapperInterfaceSliceCmder(cmd *baseCmder) ICmder {
	interfaceSliceCmd := &InterfaceSliceCmder{
		baseCmder: cmd,
	}
	v, _ := cmd.Val().([]interface{})
	interfaceSliceCmd.setVal(v)
	return interfaceSliceCmd
}

///////////////////////////////// InterfaceSliceCmder ///////////////

///////////////////////////////// ZMemberSliceCmder /////////////////
type ZMemberSliceCmder struct {
	*baseCmder
	val []ZMember
}

func (c *ZMemberSliceCmder) setVal(v []ZMember) {
	c.val = v
}

func (c *ZMemberSliceCmder) Val() []ZMember {
	return c.val
}

func (c *ZMemberSliceCmder) Result() ([]ZMember, error) {
	return c.val, c.err
}

func (c *ZMemberSliceCmder) ToString() string {
	return cmdString(c.baseCmder, c.val)
}
func wrapperZMemberSliceCmder(cmd *baseCmder) ICmder {
	zMemberSliceCmd := &ZMemberSliceCmder{
		baseCmder: cmd,
	}
	v, _ := cmd.Val().([]ZMember)
	zMemberSliceCmd.setVal(v)
	return zMemberSliceCmd
}

///////////////////////////////// ZMemberSliceCmder /////////////////

///////////////////////////////// StringStructMapCmder //////////////
type StringStructMapCmder struct {
	*baseCmder
	val map[string]struct{}
}

func (c *StringStructMapCmder) setVal(v map[string]struct{}) {
	c.val = v
}

func (c *StringStructMapCmder) Val() map[string]struct{} {
	return c.val
}

func (c *StringStructMapCmder) Result() (map[string]struct{}, error) {
	return c.val, c.err
}

func (c *StringStructMapCmder) ToString() string {
	return cmdString(c.baseCmder, c.val)
}
func wrapperStringStructMapCmder(cmd *baseCmder) ICmder {
	stringStructMapCmder := &StringStructMapCmder{
		baseCmder: cmd,
	}
	v, _ := cmd.Val().(map[string]struct{})
	stringStructMapCmder.setVal(v)
	return stringStructMapCmder
}

///////////////////////////////// StringStructMapCmder //////////////

///////////////////////////////// ScanCmder /////////////////////////
type ScanCmder struct {
	*baseCmder
	page   []string
	cursor uint64
}

func (c *ScanCmder) setVal(page []string, cursor uint64) {
	c.page = page
	c.cursor = cursor
}

func (c *ScanCmder) Val() ([]string, uint64) {
	return c.page, c.cursor
}

func (c *ScanCmder) Result() ([]string, uint64, error) {
	return c.page, c.cursor, c.err
}

func (c *ScanCmder) ToString() string {
	return cmdString(c.baseCmder, c.page)
}

// TODO result
func wrapperScanCmder(cmd *baseCmder) ICmder {
	scanCmder := &ScanCmder{
		baseCmder: cmd,
	}
	return scanCmder
}

///////////////////////////////// ScanCmder /////////////////////////
