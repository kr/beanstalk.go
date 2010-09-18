---
layout: base
title: beanstalk.go
---
<!--
Copyright 2009 The Go Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
-->

<!-- PackageName is printed as title by the top-level template -->
<p><code>import "github.com/kr/beanstalk.go.git"</code></p>
<p>
Client library for the beanstalkd protocol.
See <a href="http://kr.github.com/beanstalkd/">http://kr.github.com/beanstalkd/</a>
</p>
<p>
We are lenient about the protocol -- we accept either CR LF or just LF to
terminate server replies. We also trim white space around words in reply
lines.
</p>
<p>
This package is synchronized internally. It is safe to call any of these
functions from any goroutine at any time.
</p>
<p>
Note that, as of version 1.4.4, beanstalkd provides only 1-second
granularity on all duration values.
</p>

<p>
<h4>Package files</h4>
<span style="font-size:90%">
<a href="/src/pkg/beanstalk/beanstalk.go">beanstalk.go</a>
</span>
</p>
<h2 id="Constants">Constants</h2>
<p>
For use in parameters that measure duration (in microseconds). Not really
infinite; merely large. About 63 years.
</p>

<pre>const Forever = 2000000000000000 <span class="comment">// µs</span>
</pre>
<h2 id="Variables">Variables</h2>
<p>
Reasons for an invalid tube name.
</p>

<pre>var (
    NameTooLong = os.NewError(&#34;name too long&#34;)
    IllegalChar = os.NewError(&#34;name contains illegal char&#34;)
)</pre>
<p>
Error responses from the server.
</p>

<pre>var (
    OutOfMemory   = os.NewError(&#34;Server Out of Memory&#34;)
    InternalError = os.NewError(&#34;Server Internal Error&#34;)
    Draining      = os.NewError(&#34;Server Draining&#34;)
    Buried        = os.NewError(&#34;Buried&#34;)
    JobTooBig     = os.NewError(&#34;Job Too Big&#34;)
    TimedOut      = os.NewError(&#34;Reserve Timed Out&#34;)
    NotFound      = os.NewError(&#34;Job or Tube Not Found&#34;)
    NotIgnored    = os.NewError(&#34;Tube Not Ignored&#34;)
)</pre>
<p>
The server sent a bad reply. For example: unknown or inappropriate
response, wrong number of terms, or invalid format.
</p>

<pre>var BadReply = os.NewError(&#34;Bad Reply from Server&#34;)</pre>
<h2 id="Conn">type <a href="/src/pkg/beanstalk/beanstalk.go#L31">Conn</a></h2>
<p>
A connection to beanstalkd. Provides methods that operate outside of any
tube. This type also embeds Tube and TubeSet, which is convenient if you
rarely change tubes.
</p>

<p><pre>type Conn struct {
    Name string
    *Tube
    *TubeSet
    <span class="comment">// contains unexported fields</span>
}</pre></p>
<h3 id="Conn.Dial">func <a href="/src/pkg/beanstalk/beanstalk.go#L419">Dial</a></h3>
<p><code>func Dial(addr string) (*Conn, os.Error)</code></p>
<p>
Dial the beanstalkd server at remote address addr.
</p>

<h3 id="Conn.ListTubes">func (*Conn) <a href="/src/pkg/beanstalk/beanstalk.go#L620">ListTubes</a></h3>
<p><code>func (c *Conn) ListTubes() ([]string, os.Error)</code></p>

<h3 id="Conn.Peek">func (*Conn) <a href="/src/pkg/beanstalk/beanstalk.go#L612">Peek</a></h3>
<p><code>func (c *Conn) Peek(id uint64) (*Job, os.Error)</code></p>
<p>
Get a copy of the specified job.
</p>

<h3 id="Conn.Stats">func (*Conn) <a href="/src/pkg/beanstalk/beanstalk.go#L616">Stats</a></h3>
<p><code>func (c *Conn) Stats() (map[string]string, os.Error)</code></p>

<h2 id="Error">type <a href="/src/pkg/beanstalk/beanstalk.go#L60">Error</a></h2>
<p>
Implements os.Error
</p>

<p><pre>type Error struct {
    ConnName string
    Cmd      string
    Reply    string
    Error    os.Error
}</pre></p>
<h3 id="Error.String">func (Error) <a href="/src/pkg/beanstalk/beanstalk.go#L88">String</a></h3>
<p><code>func (e Error) String() string</code></p>

<h2 id="Job">type <a href="/src/pkg/beanstalk/beanstalk.go#L38">Job</a></h2>

<p><pre>type Job struct {
    Id   uint64
    Body string
    <span class="comment">// contains unexported fields</span>
}</pre></p>
<h3 id="Job.Bury">func (*Job) <a href="/src/pkg/beanstalk/beanstalk.go#L714">Bury</a></h3>
<p><code>func (j *Job) Bury(pri uint32) os.Error</code></p>
<p>
Bury job j and change its priority to pri.
</p>

<h3 id="Job.Delete">func (*Job) <a href="/src/pkg/beanstalk/beanstalk.go#L704">Delete</a></h3>
<p><code>func (j *Job) Delete() os.Error</code></p>
<p>
Delete job j.
</p>

<h3 id="Job.Release">func (*Job) <a href="/src/pkg/beanstalk/beanstalk.go#L719">Release</a></h3>
<p><code>func (j *Job) Release(pri uint32, µsDelay uint64) os.Error</code></p>
<p>
Release job j, changing its priority to pri and its delay to delay.
</p>

<h3 id="Job.Stats">func (*Job) <a href="/src/pkg/beanstalk/beanstalk.go#L725">Stats</a></h3>
<p><code>func (j *Job) Stats() (map[string]string, os.Error)</code></p>
<p>
Get statistics on job j.
</p>

<h3 id="Job.Touch">func (*Job) <a href="/src/pkg/beanstalk/beanstalk.go#L709">Touch</a></h3>
<p><code>func (j *Job) Touch() os.Error</code></p>
<p>
Touch job j.
</p>

<h2 id="Tube">type <a href="/src/pkg/beanstalk/beanstalk.go#L46">Tube</a></h2>
<p>
Represents a single tube. Provides methods that operate on one tube,
especially Put.
</p>

<p><pre>type Tube struct {
    Name string
    <span class="comment">// contains unexported fields</span>
}</pre></p>
<h3 id="Tube.NewTube">func <a href="/src/pkg/beanstalk/beanstalk.go#L629">NewTube</a></h3>
<p><code>func NewTube(c *Conn, name string) (*Tube, os.Error)</code></p>
<p>
Returns an error if the tube name is invalid.
</p>

<h3 id="Tube.Kick">func (*Tube) <a href="/src/pkg/beanstalk/beanstalk.go#L691">Kick</a></h3>
<p><code>func (t *Tube) Kick(n uint64) (uint64, os.Error)</code></p>
<p>
Kick up to n jobs in tube t.
</p>

<h3 id="Tube.Pause">func (*Tube) <a href="/src/pkg/beanstalk/beanstalk.go#L696">Pause</a></h3>
<p><code>func (t *Tube) Pause(µs uint64) os.Error</code></p>
<p>
Pause tube t for µs microseconds.
</p>

<h3 id="Tube.PeekBuried">func (*Tube) <a href="/src/pkg/beanstalk/beanstalk.go#L679">PeekBuried</a></h3>
<p><code>func (t *Tube) PeekBuried() (*Job, os.Error)</code></p>
<p>
Get a copy of a buried job in this tube, if any.
</p>

<h3 id="Tube.PeekDelayed">func (*Tube) <a href="/src/pkg/beanstalk/beanstalk.go#L674">PeekDelayed</a></h3>
<p><code>func (t *Tube) PeekDelayed() (*Job, os.Error)</code></p>
<p>
Get a copy of the next delayed job in this tube, if any.
</p>

<h3 id="Tube.PeekReady">func (*Tube) <a href="/src/pkg/beanstalk/beanstalk.go#L669">PeekReady</a></h3>
<p><code>func (t *Tube) PeekReady() (*Job, os.Error)</code></p>
<p>
Get a copy of the next ready job in this tube, if any.
</p>

<h3 id="Tube.Put">func (*Tube) <a href="/src/pkg/beanstalk/beanstalk.go#L471">Put</a></h3>
<p><code>func (t *Tube) Put(body string, pri uint32, µsDelay, µsTTR uint64) (id uint64, err os.Error)</code></p>
<p>
Put a job into the queue and return its id.
</p>
<p>
If an error occured, err will be non-nil. For some errors, Put will also
return a valid job id, so you must check both values.
</p>
<p>
The delay and ttr are measured in microseconds.
</p>

<h3 id="Tube.Stats">func (*Tube) <a href="/src/pkg/beanstalk/beanstalk.go#L684">Stats</a></h3>
<p><code>func (t *Tube) Stats() (map[string]string, os.Error)</code></p>
<p>
Get statistics on tube t.
</p>

<h2 id="TubeError">type <a href="/src/pkg/beanstalk/beanstalk.go#L67">TubeError</a></h2>

<p><pre>type TubeError struct {
    TubeName string
    Error    os.Error
}</pre></p>
<h3 id="TubeError.String">func (TubeError) <a href="/src/pkg/beanstalk/beanstalk.go#L92">String</a></h3>
<p><code>func (e TubeError) String() string</code></p>

<h2 id="TubeSet">type <a href="/src/pkg/beanstalk/beanstalk.go#L53">TubeSet</a></h2>
<p>
Represents a set of tubes. Provides methods that operate on several tubes at
once, especially Reserve.
</p>

<p><pre>type TubeSet struct {
    Names []string
    <span class="comment">// contains unexported fields</span>
}</pre></p>
<h3 id="TubeSet.NewTubeSet">func <a href="/src/pkg/beanstalk/beanstalk.go#L640">NewTubeSet</a></h3>
<p><code>func NewTubeSet(c *Conn, names []string) (*TubeSet, os.Error)</code></p>
<p>
Returns an error if any of the tube names are invalid.
</p>

<h3 id="TubeSet.Reserve">func (*TubeSet) <a href="/src/pkg/beanstalk/beanstalk.go#L653">Reserve</a></h3>
<p><code>func (t *TubeSet) Reserve() (*Job, os.Error)</code></p>
<p>
Reserve a job from any one of the tubes in t.
</p>

