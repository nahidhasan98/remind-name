function showToast(type, message) {
    console.log(message);
    $(".toast-body").html(message.replace(/\n/g, "<br>"));
    if (type == "success") $(".toast").removeClass("toastFail").toast('show');
    else $(".toast").addClass("toastFail").toast('show');
}

function toggleLoadingIcon(todo, form) {
    if (todo == "show") {
        // Disable the submit button and show the loading gif
        $("#" + form + " .submitBtn").prop("disabled", true);
        $("#" + form + " .loadingIcon").show();
    } else {
        // Re-enable the submit button and hide the loading gif
        $("#" + form + " .submitBtn").prop("disabled", false);
        $("#" + form + " .loadingIcon").hide();
    }
}

function copyBtnEvent() {
    $('#copyBtn').on("click", () => {
        if (navigator.clipboard) {
            navigator.clipboard.writeText($('#token').text().trim())
                .then(() => showToast("success", "Text copied to clipboard!"))
                .catch(err => alert('Failed to copy text: ', err));
        } else {
            showToast("error", "Copy Text button won't work because<br>Clipboard API not supported on this device or browser.")
        }
    });
}

function defaultScheduleBtnEvent() {
    $("#defaultSchedule").on("click", (event) => {
        $("#defaultSchedule").addClass("active");
        $("#customSchedule").removeClass("active");
        $("#defaultUncheck, #customCheck").hide();
        $("#defaultCheck, #customUncheck").show();
        $("#defaultSection").slideDown();
        $("#customSection").slideUp();
        // Set the hidden input value to "default"
        $("#scheduleType").val("default");
    });
}

function customScheduleBtnEvent() {
    $("#customSchedule").on("click", (event) => {
        $("#defaultSchedule").removeClass("active");
        $("#customSchedule").addClass("active");
        $("#customUncheck, #defaultCheck").hide();
        $("#customCheck, #defaultUncheck").show();
        $("#defaultSection").slideUp();
        $("#customSection").slideDown();
        // Set the hidden input value to "custom"
        $("#scheduleType").val("custom");
    });
}

function toggleAMPM(container) {
    $(container + " span").on("click", (event) => {
        $(event.target).css({
            "color": "#fff",
            "background": "#355e97",
            "border": "1px solid #355e97"
        });

        $(container + " span").not(event.target).css({
            "color": "#122033",
            "background": "none",
            "border": "1px solid #9eb0c9"
        });

        // Set the hidden input value based on selected AM/PM
        const selectedValue = $(event.target).hasClass('am') ? 'am' : 'pm';
        $(container + " input").val(selectedValue);
    });
}

function onClickEvent() {
    copyBtnEvent();
    defaultScheduleBtnEvent();
    customScheduleBtnEvent();
    toggleAMPM("#fromAMPM");
    toggleAMPM("#toAMPM");
}

function platformSelectEvent(instance) {
    $("#platform").on("change", (event) => {
        instance.update();
        let value = $(event.target).val()
        let placeholder = "";
        if (value == "Telegram") placeholder = "@RemindNameBot";

        $("#usernameLabel").html("Your " + value + " Username/ID");
        $("#username").attr("placeholder", placeholder);
    });
}

function onChangeEvent(instance) {
    platformSelectEvent(instance);
}

function scheduleHoverEffect(elem) {
    $(elem).hover(
        () => {
            $(elem).css("background", "#b3d0ff");
        }, () => {
            $(elem).css("background", "#cae0ff");
        }
    );
}

function handleHoverEffect() {
    scheduleHoverEffect("#defaultSchedule");
    scheduleHoverEffect("#customSchedule");
}

function setTimeZones() {
    $("#timezone").val(Intl.DateTimeFormat().resolvedOptions().timeZone);
}

function isValidFormData(form) {
    if (form == "subscriptionForm") {
        // taking care platform select box
        let val = $('#platform').val().trim();
        if (val == "Select an option" || val.length == 0) {
            showToast("error", "Please choose a Platform!");
            return false;
        }
        // taking care of username
        if ($('#username').val().trim().length == 0) {
            showToast("error", "Please enter your " + val + " username!");
            return false;
        }

        return true
    } else if (form == "feedbackForm") {
        if ($('#name').val().trim().length == 0) {
            showToast("error", "Please enter your name!");
            return false;
        }

        if ($('#email').val().trim().length == 0) {
            showToast("error", "Please enter your email!");
            return false;
        }

        if ($('#feedbackText').val().trim().length == 0) {
            showToast("error", "Please enter your feedback!");
            return false;
        }

        return true
    }

    return false;
}

function subscriptionFormHandling() {
    let form = "subscriptionForm";
    $("#" + form).on("submit", (event) => {
        event.preventDefault();
        toggleLoadingIcon("show", form);

        if (!isValidFormData(form)) {
            toggleLoadingIcon("hide", form);
            return false;
        }

        let formData = $('form').serialize();
        // console.log(formData);

        // sending ajax post request
        let request = $.ajax({
            type: "POST",
            url: "/subscription",
            data: formData,
        });

        request.done(function (response) {
            // console.log(response);
            showToast("success", response.Message);

            if (response.Status == 0) {
                $("#subscriptionStatus span").text("Not Verified");
                $("#subscriptionText").text(response.Message);
                $(".platformName").text(response.Platform.Name);
                $("#bot_name").text(response.Platform.BotName);
                $("#bot_username").text(response.Platform.BotUsername);
                $("#bot_link").html("<a href='" + response.Platform.BotLink + "'>" + response.Platform.BotLink + "</a>");
                $("#token span").text(response.Token);
                $("#botSection").slideDown();
            } else if (response.Status == 1) {
                $("#subscriptionStatus span").text("Verified");
                $("#botSection").slideUp();
            }
        });

        request.fail(function (response) {
            // console.log(response);
            showToast("error", response.responseJSON.error);
        });

        request.always(function () {
            toggleLoadingIcon("hide", form);
        });

        return false;
    });
}

function feedbackFormHandling() {
    let form = "feedbackForm";
    $("#" + form).on("submit", (event) => {
        event.preventDefault();
        toggleLoadingIcon("show", form);

        if (!isValidFormData(form)) {
            toggleLoadingIcon("hide", form);
            return false;
        }

        let formData = $('form').serialize();
        // console.log(formData);

        // sending ajax post request
        let request = $.ajax({
            type: "POST",
            url: "/feedback",
            data: formData,
        });

        request.done(function (response) {
            // console.log(response);
            showToast("success", response.Message);
            $("#feedbackModal").modal("hide");
        });

        request.fail(function (response) {
            // console.log(response);
            showToast("error", response.responseJSON.error);
        });

        request.always(function () {
            toggleLoadingIcon("hide", form);
        });

        return false;
    });
}

function formSubmitHandling() {
    subscriptionFormHandling();
    feedbackFormHandling();
}

$(document).ready(function () {
    let instance = NiceSelect.bind(document.getElementById("platform"));
    $('[data-bs-toggle="tooltip"]').tooltip()

    onClickEvent();
    onChangeEvent(instance);
    handleHoverEffect();
    setTimeZones();
    formSubmitHandling();
});