$('.digit-group').find('input').each(function () {
    $(this).attr('maxlength', 1);
    $(this).on('keyup', function (e) {
        var parent = $($(this).parent());

        if (e.keyCode === 8 || e.keyCode === 37) {
            var prev = parent.find('input#' + $(this).data('previous'));

            if (prev.length) {
                $(prev).select();
            }
        } else if ((e.keyCode >= 48 && e.keyCode <= 57) || (e.keyCode >= 65 && e.keyCode <= 90) || (e.keyCode >= 96 && e.keyCode <= 105) || e.keyCode === 39) {
            var next = parent.find('input#' + $(this).data('next'));

            if (next.length) {
                $(next).select();
            } else {
                if (parent.data('autosubmit')) {
                    parent.submit();
                }
            }
        }
    });
});

let error = document.getElementById("error");

function getOtp() {
    return document.getElementById("digit-1").value + document.getElementById("digit-2").value + document.getElementById("digit-3").value +
        document.getElementById("digit-4").value + document.getElementById("digit-5").value + document.getElementById("digit-6").value;
}

async function sendOtp() {
    let credentials = localStorage.getItem('credentials');
    let otp = getOtp();

    let response = await fetch('/otp/check/' + otp, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: credentials
    });

    let result = await response.text();
    let contentType = response.headers.get("Content-Type");

    if (!contentType.includes("text/html")) {
        const respJson = JSON.parse(result);
        error.innerHTML = respJson.msg;
        error.style.display = "block";
    } else {
        error.style.display = "none";
        document.write(result);
        document.close();
    }
}