"use strict"

let error = document.getElementById("error")

function togglePasswordVisibility() {
    let toggle = document.getElementById("togglePassword");
    let password = document.getElementById("password")

    if (toggle.checked) {
        password.type = "text";
    } else {
        password.type = "password";
    }
}

function validateEmail(email) {
    var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(String(email).toLowerCase());
}

function getEmailAndPassword() {
    let email = document.getElementById("email").value;
    let password = document.getElementById("password").value;

    //check if its a valid email
    if (!validateEmail(email)) {
        error.style.display = "block";
        error.innerHTML = "Enter a valid email";
        return false;
    }

    if (password.length <= 0) {
        error.style.display = "block";
        error.innerHTML = "Enter a password";
        return false;
    }

    return { email: email, password: password };
}

async function login() {
    let data = getEmailAndPassword();

    if (!data) {
        return;
    }

    let result = await fetch('/users/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    });

    result = await result.json();

    if (result.code == 200) {
        error.style.display = "none";
        error.innerHTML = "";
        window.location.href = "/otp";
        return
    }

    error.style.display = "block";
    error.innerHTML = result.msg;
}