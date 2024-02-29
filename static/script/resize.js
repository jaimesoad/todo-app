/**
 * @param {string} div
 * @param {number} minHeight   
 * @param {number} minWidth 
 */
function makeResizableDiv(div, minHeight, minWidth) {
    const element = document.querySelector(div);

    for (let hor of ["left", "right"]){
        for (let ver of ["top", "bottom"]) {
            const newChild = document.createElement("div")
            newChild.className = `corner ${ver}-${hor}`

            element.appendChild(newChild)
        }
    }

    const resizers = document.querySelectorAll(div + ' .corner')
    
    let originalWidth = 0;
    let originalHeight = 0;
    let originalX = 0;
    let originalY = 0;
    let originalMouseX = 0;
    let originalMouseY = 0;
    for (let i = 0; i < resizers.length; i++) {
        const currentResizer = resizers[i];
        currentResizer.addEventListener('mousedown', /** @param {MouseEvent} e */ (e) => {
            e.preventDefault()
            originalWidth = parseFloat(getComputedStyle(element, null).getPropertyValue('width').replace('px', ''));
            originalHeight = parseFloat(getComputedStyle(element, null).getPropertyValue('height').replace('px', ''));
            originalX = element.getBoundingClientRect().left;
            originalY = element.getBoundingClientRect().top;
            originalMouseX = e.pageX;
            originalMouseY = e.pageY;
            window.addEventListener('mousemove', resize)
            window.addEventListener('mouseup', stopResize)
        })

        /** @param {MouseEvent} e */
        const resize = (e) => {
            if (currentResizer.classList.contains('bottom-right')) {
                const width = originalWidth + (e.pageX - originalMouseX);
                const height = originalHeight + (e.pageY - originalMouseY)
                if (width > minWidth) {
                    element.style.width = width + 'px'
                }
                if (height > minHeight) {
                    element.style.height = height + 'px'
                }
            }
            else if (currentResizer.classList.contains('bottom-left')) {
                const height = originalHeight + (e.pageY - originalMouseY)
                const width = originalWidth - (e.pageX - originalMouseX)
                if (height > minHeight) {
                    element.style.height = height + 'px'
                }
                if (width > minWidth) {
                    element.style.width = width + 'px'
                    element.style.left = originalX + (e.pageX - originalMouseX) + 'px'
                }
            }
            else if (currentResizer.classList.contains('top-right')) {
                const width = originalWidth + (e.pageX - originalMouseX)
                const height = originalHeight - (e.pageY - originalMouseY)
                if (width > minWidth) {
                    element.style.width = width + 'px'
                }
                if (height > minHeight) {
                    element.style.height = height + 'px'
                    element.style.top = originalY + (e.pageY - originalMouseY) + 'px'
                }
            }
            else {
                const width = originalWidth - (e.pageX - originalMouseX)
                const height = originalHeight - (e.pageY - originalMouseY)
                if (width > minWidth) {
                    element.style.width = width + 'px'
                    element.style.left = originalX + (e.pageX - originalMouseX) + 'px'
                }
                if (height > minHeight) {
                    element.style.height = height + 'px'
                    element.style.top = originalY + (e.pageY - originalMouseY) + 'px'
                }
            }
            localStorage.setItem("loginSize", JSON.stringify({
                height: element.style.height,
                width: element.style.width
            }))
        }

        function stopResize() {
            window.removeEventListener('mousemove', resize)
        }
    }
}

makeResizableDiv('#window', 240, 480)