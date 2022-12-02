/* dropdown-menu */

document.addEventListener("click", e => {
    const isDropdownBtn = e.target.matches("[data-dropdown-btn]");
    if(!isDropdownBtn && e.target.closest("[data-dropdown]") != null){
        return;
    }

    let currentDropdown
    if(isDropdownBtn) {
        currentDropdown = e.target.closest("[data-dropdown]");
        currentDropdown.classList.toggle("menu-open");
    }

    document.querySelectorAll("[data-dropdown].menu-open").forEach(dropdown => {
        if(dropdown === currentDropdown){
            return;
        }
        dropdown.classList.remove("menu-open");
    })

    
})

/* responsive nav bar */

var responsive = document.querySelector(".responsive");

// if(responsive) {
//     var navbtn = responsive.querySelector("[data-navbar-btn]");
//     var navContent = responsive.querySelector(".responsive-navbar");

//     navbtn.addEventListener("click", () => {
//         console.log("test");
//         navContent.classList.add("active");
//     })
// }

document.addEventListener("click", e => {
    const isNavBtn = e.target.matches("[data-navbar-btn]");
    if(!isNavBtn && e.target.closest("[data-navbar]") != null){
        return;
    }

    let currentNav
    if(isNavBtn) {
        currentNav = e.target.closest("[data-navbar]");
        currentNav.classList.toggle("active");

        let clsNavBtn = currentNav.querySelector(".close-btn");
        clsNavBtn.addEventListener("click", () => {
            currentNav.classList.remove("active");
        })
    }

    document.querySelectorAll("[data-navbar].active").forEach(NavBar => {
        if(NavBar === currentNav){
            return;
        }
        NavBar.classList.remove("active");
    })
    
    
    
})

/* collapsible */

var coll = document.querySelectorAll("[data-collapsible-btn]");
var i;

for (i = 0; i < coll.length; i++) {
    coll[i].addEventListener("click", function() {
        this.classList.toggle("active");
        var content = this.nextElementSibling;
        
        if (content.style.maxHeight){
        content.style.maxHeight = null;
        
        } else {
        content.style.maxHeight = content.scrollHeight + "px";
        
        } 
    });
};

