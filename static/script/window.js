const body = document.querySelector("body")
const login = document.getElementById("window")
const topbar = document.getElementById("topbar")

/** @type {{height: string, width: string} | null}  */
const size = JSON.parse(localStorage.getItem("loginSize"))

/** @type {{top: string, left: string} | null}  */
const position = JSON.parse(localStorage.getItem("loginPosition"))

if (size !== null) {
    login.style.height = size.height
    login.style.width = size.width

} else if (window.innerWidth < 500) {
    login.style.width = `${window.innerWidth - 20}px`

}
if (position !== null) {
    login.style.top = position.top
    login.style.left = position.left

} else {
    // places the window right in the center of the screen.
    login.style.top = (window.innerHeight - login.offsetHeight) / 2 + "px"
    login.style.left = (window.innerWidth - login.offsetWidth) / 2 + "px"
}

let x = 0, y = 0
let mousedown = false

topbar.addEventListener("mousedown", (e) => {
    mousedown = true

    x = login.offsetLeft - e.clientX
    y = login.offsetTop - e.clientY

    e.preventDefault()
}, false)

topbar.addEventListener("mouseup", () => mousedown = false, false)

body.addEventListener("mousemove", (e) => {
    if (mousedown) {
        login.style.left = `${e.clientX + x}px`
        login.style.top = `${e.clientY + y}px`

        localStorage.setItem("loginPosition", JSON.stringify({
            top: login.style.top,
            left: login.style.left
        }))
    }
}, false)