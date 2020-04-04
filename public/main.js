function ping() {
    fetch('/api/v1/ping').then(res => res.json()).then(data => {
        document.getElementById('version').innerText = data.go_version;
    });
}

let msgbox = document.getElementById('outputMsg')
let ssabox = document.getElementById('ssa')
let lastFuncName, lastCode;
function build() {
    let funcname = document.getElementById('funcname').value;
    let code = document.getElementById('code').value;

    // several early checks
    if (funcname === lastFuncName && code == lastCode) {
        console.log('no changes, do not submit')
        return
    }
    if (!code.includes('func '+funcname)) {
        ssabox.src = ''
        msgbox.innerText = 'GOFUNCNAME does not exist in your code.'
        return
    }

    msgbox.innerText = ''
    lastFuncName = funcname
    lastCode = code

    fetch('/api/v1/buildssa', {
        method: 'POST',
        headers: {'Content-Type': 'application/json',},
        body: JSON.stringify({
            'funcname': funcname,
            'code': code,
        }),
    })
    .then((response) => {
        return new Promise((resolve, reject) => {
            // will resolve or reject depending on status, will pass both "status" and "data" in either case
            let func;
            response.status < 400 ? func = resolve : func = reject;
            response.json().then(data => func({'status': response.status, 'data': data}));
        });
    })
    .then(res => {
        ssabox.src = `./buildbox/${res.data.build_id}/ssa.html`
        msgbox.innerText = ''
    })
    .catch(res => {
        ssabox.src = ''
        msgbox.innerText = res.data.msg
    });
}

// inject ssa style
let ssablock = document.getElementById('ssa')
ssablock.addEventListener('load', () => {
    // wait until ssa is loaded
    let $head = $("iframe").contents().find("head");
    $head.append($("<link/>", { rel: "stylesheet", href: "/scrollbar.css", type: "text/css"}));
})

// inject version info
ping();

// listen build action
let buildssa = document.getElementById('build')
buildssa.addEventListener('click', () => {
    build()
})

// listen about action
let aboutinfo = $('#about');
aboutinfo.click(function(e) {
    if ($(e.target).is('a')) {
        return;
    }
    aboutinfo.hide();
});
$('#aboutbtn').click(function() {
    if (aboutinfo.is(':visible')) {
        aboutinfo.hide();
        return;
    }
    aboutinfo.show();
})

$('#code').linedtextarea();
$('#code').keydown(function(event){
    if (event.keyCode == 9) {
        event.preventDefault();
        var start = this.selectionStart;
        var end = this.selectionEnd;
        // set textarea value to: text before caret + tab + text after caret
        $(this).val($(this).val().substring(0, start)
                    + "\t"
                    + $(this).val().substring(end));
        // put caret at right position again
        this.selectionStart =
        this.selectionEnd = start + 1;
    }
});

// TODO: dragable scroll
// let wholePage = document.querySelector('body');
// let el = document.querySelector("#ssa").contentDocument.querySelector('body');
// let x = 0, y = 0, t = 0, l = 0;
// console.log(el);

// let draggingFunction = (e) => {
//     wholePage.addEventListener('mouseup', () => {
//         wholePage.removeEventListener("mousemove", draggingFunction);
//     });

//     el.scrollLeft = l - e.pageX + x;
//     el.scrollTop = t - e.pageY + y;
// };

// wholePage.addEventListener('click', (e) => {
//     console.log("xxx")

//     y = e.pageY;
//     x = e.pageX;
//     t = el.scrollTop;
//     l = el.scrollLeft;

//     el.addEventListener('mousemove', draggingFunction);
// });