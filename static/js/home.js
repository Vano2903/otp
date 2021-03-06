let user = JSON.parse(localStorage.getItem('credentials'));

async function loadPage() {
    document.getElementById("user-welcome").value += user.email;

    let response = await fetch("/users/pfp/" + user.email, {
        method: 'GET',
    });
    let resp = await response.json();
    if (resp.code == 200) {
        document.getElementById("pfp").src = resp.pfpUrl;
    } else {
        log.Error("Error fetching: " + resp.msg)
    }
}

function showPfpUpdateForm() {
    if (document.getElementById("updatePfpForm").style.display == "none") {
        document.getElementById("updatePfpForm").style.display = "block";
    } else {
        document.getElementById("updatePfpForm").style.display = "none";
    }
}

async function uploadFile() {
    const formData = new FormData();
    let file = document.getElementById("newPfpUrl").files[0]
    formData.append("document", file);
    const response = await fetch("/upload/file", {
        method: 'POST',
        body: formData
    });
    let resp = await response.json();
    console.log(resp)
    return {
        id: resp.fileID, name: file.name
    }
}

async function uploadInfo() {
    let fileData = await uploadFile()
    console.log(fileData)
    let body = {
        email: user.email,
        password: user.password,
        id: fileData.id,
    }
    console.log(body)

    const response = await fetch("/upload/bind/" + fileData.id, {
        method: 'POST',
        body: JSON.stringify(body)
    });
    let resp = await response.json();

    if (resp.code == 202) {
        document.getElementById("error").style.display = "none";
        document.getElementById("pfp").src = resp.pfpUrl;
    } else {
        document.getElementById("error").style.display = "block";
        document.getElementById("error").innterHTML = resp.msg;
    }

    console.log(resp)
}


// async function updatePassword() {
//     currentPass = document.getElementById("currentPassword").value;
//     newPass = document.getElementById("newPassword").value;
//     confirmNewPass = document.getElementById("confirmNewPassword").value;
//     if (currentPass != user.password) {
//         $("#errCurrentPass").innerHTML = "The password is different from the current one";
//         return
//     }
//     if (newPass == "") {
//         $("#errNewPassword").innerHTML = "New password can't be empty";
//         return
//     }
//     if (newPass == currentPass) {
//         $("#errNewPassword").innerHTML = "The new password can't be equal to the current one";
//         return
//     }
//     if (newPass != confirmNewPass) {
//         $("#errConfirmPass").innerHTML = "The two passwords do not match";
//         return
//     }
//     const res = await fetch('/users/customization/?password=' + newPass, {
//         method: "POST",
//         headers: {
//             'Content-Type': 'application/json'
//         },
//         body: JSON.stringify(user)
//     });
//     const resp = await res.json();
//     if (res.status != 200) {
//         alert("qualcosa ?? andato storto, riprova")
//         return
//     }
//     alert("modifica andata a buon fine, verr?? richiesto di fare login nuovamente")
//     //logout()
// }