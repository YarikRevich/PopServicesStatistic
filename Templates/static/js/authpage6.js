const timer = ms => new Promise(res => setTimeout(res, ms))

async function rotation() { // We need to wrap the loop into an async function for this to work
    var elem = document.getElementById("auth-part");
    elem.removeAttribute("onclick");
    for (let i = 0; i <= 360; i++) {
        elem.style.height = i * 0.9 + "px";
        elem.style.width = i * 0.9 + "px";
        elem.style.borderRadius = i / 4.5 + "px";
        elem.style.transform = "rotate(" + i + "deg)";
        await timer(0.5); // then the created Promise can be awaited
    }
    for (let i = 0; i <= 360; i++) {
        document.getElementById("auth-child").style.opacity = i / 360;
        await timer(0.5);
    }
}

const request = new XMLHttpRequest()



request.open("POST", "http://localhost:8000/usability")

request.setRequestHeader("Content-Type", "application/json")


request.addEventListener("readystatechange", () => {

    if (request.readyState === 4 && request.status === 200) {
        console.log(request.responseText)
    }
})

request.send()