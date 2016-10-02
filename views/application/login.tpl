
<div class="dialog">
    <div class="dialog-header">
        <h1>Login</h1>
    </div>
    <div class="dialog-body">
        <form action="{{ .Prefix }}/login" method="POST" autocomplete="off" accept-charset="utf-8">
            <input class="text" type="text" name="name" autofocus="autofocus" tabindex="1" placeholder="name" required />
            <input class="text" type="password" name="password" tabindex="2" placeholder="password" required />
            <div class="foot">
              <div class="container">
                <a id="license-link" href="{{ .Prefix }}/license" title="This is free software. Click this link for more information">Source code &amp; license</a>
              </div>
              <input class="btn yellow-btn" type="submit" tabindex="3" value="login" />
              <div class="clearout"></div>
            </div>
            <div class="clearout"></div>
        </form>
    </div>
</div>
