
<div class="topics-wrapper">


<div id="list">
    <ul>
        {{range .Topics}}
            <li><a href="/topics/{{.Id}}">{{.Name}}</a></li>
        {{end}}
        <li class="hidden">
            <form action="/topics" method="POST" autocomplete="off" accept-charset="utf-8">
                <input class="text" type="text" name="name" value="" autocomplete="off" tabindex="1" placeholder="Nom" />
                <input class="btn yellow-btn" type="submit" value="Create" />
                <div class="clearout"></div>
            </form>
        </li>
    </ul>
</div>

<div id="contents">
    <div class="contents-body">
        {{noescape .Rendered}}
    </div>

    <div class="contents-edit">
        <form action="/topics/{{.Current.Id}}" method="POST" autocomplete="off" accept-charset="utf-8">
            <textarea name="contents" autocomplete="off">{{.Current.Contents}}</textarea>
            <input class="btn yellow-btn" type="submit" value="Update" />
        </form>
    </div>
</div>


</div>

