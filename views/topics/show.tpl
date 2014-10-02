
<div class="topics-wrapper">

{{$c := .Current}}

<div id="list">
    <h3>topics</h3>
    <ul>
        {{range .Topics}}
            {{if eq .Id $c.Id}}
                <li class="selected"><a href="/topics/{{.Id}}">{{.Name}}</a></li>
            {{else}}
                <li><a href="/topics/{{.Id}}">{{.Name}}</a></li>
            {{end}}
        {{end}}

        <li class="create">
            <button class="btn grey-btn">create</button>

            <form action="/topics" method="POST" autocomplete="off" accept-charset="utf-8">
                <input class="tiny-text" type="text" name="name" value="" autocomplete="off" tabindex="1" placeholder="name" />
                <input class="btn grey-btn" type="submit" value="create" />
            </form>
        </li>

        <li class="logout">
            <form action="/logout" method="POST">
                <input class="btn grey-btn" type="submit" value="logout" />
            </form>
        </li>
    </ul>
</div>

<div id="contents">
{{if $c.Id}}
    <div class="header">
        <h2>{{$c.Name}}</h2>
        <ul class="options">
            <li><a href="#" class="rename" title="Rename this topic">rename</a></li>
            <li><a href="#" class="edit" title="Edit the contents of this topic">edit</a></li>
            <li><a href="#" class="delete" title="Delete this topic">delete</a></li>
            <li class="confirmation"><span>are you sure? <a href="#" class="yes">yes</a> / <a href="#" class="no">no</a></span></li>
        </ul>

        <form action="/topics/{{.Current.Id}}/delete" class="delete-form" method="POST">
        </form>

        <form action="/topics/{{.Current.Id}}" class="rename-form" method="POST" autocomplete="off" accept-charset="utf-8">
            <input class="tiny-text" type="text" name="name" value="" autocomplete="off" tabindex="1" placeholder="name" />
            <input class="btn grey-btn" type="submit" value="rename" />
        </form>
    </div>

    <div class="body">
        <div class="contents-body">
            {{noescape .Rendered}}
        </div>

        <div class="contents-edit">
            <form action="/topics/{{.Current.Id}}" method="POST" autocomplete="off" accept-charset="utf-8">
                <textarea name="contents" autocomplete="off" spellcheck="off">{{.Current.Contents}}</textarea>
                <div class="buttons">
                    <input class="btn green-btn" type="submit" value="update" />
                    <button class="btn red-btn cancel-btn">cancel</button>
                </div>
            </form>
        </div>
    </div>
{{end}}
</div>


</div>

