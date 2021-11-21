"use strict";

let error = document.getElementById("error")
let success = document.getElementById("success")


function validateEmail(email) {
    var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(String(email).toLowerCase());
}

function checkConfirmPassword() {
    var password = document.getElementById("password").value;
    var confirmPassword = document.getElementById("confirmPassword").value;
    if (password != confirmPassword) {
        document.getElementById("error").style.display = "block";
        document.getElementById("error").innerHTML = "Passwords do not match";
        return false;
    } else {
        document.getElementById("error").style.display = "none";
        document.getElementById("error").innerHTML = "";
        return true;
    }
}



async function register() {
    var email = document.getElementById("email").value;
    if (!validateEmail(email)) {
        document.getElementById("error").style.display = "block";
        document.getElementById("error").innerHTML = "Invalid Email Address";
        return
    }
    if (!checkConfirmPassword()) {
        document.getElementById("error").style.display = "block";
        document.getElementById("error").innerHTML = "The Passwords do not match";
        return;
    }
    var password = document.getElementById("password").value;

    var response = await fetch("/users/singup", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ email: email, password: password })
    });
    var data = await response.json();
    if (data.code == 202) {
        error.style.display = "none";
        success.style.display = "block";
        success.innerHTML = "Confirm the registration via email";
    } else {
        success.style.display = "none";
        error.style.display = "block";
        error.innerHTML = data.msg;
    }
}