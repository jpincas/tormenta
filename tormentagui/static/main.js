const saveEntity = (entity, id) => {
    const data = document.getElementById("entity-json").innerHTML;
    postData(`/api/${entity}/${id}`, data)
        .then(response => alertSaveResult(response));
}

const postData = (url = ``, data = "") => {
    return fetch(url, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json; charset=utf-8",

        },
        body: data,
    })
};

const alertSaveResult = response => {
    if (response.ok) {
        alertSaveSuccess();
    } else {
        response.json().then(
            error => alertSaveFail(error.errorMessage)
        )

    }
}

const alertSaveSuccess = () => {
    document.getElementById("save-result").classList.add("success");
    document.getElementById("save-result").innerHTML = `<i class="fas fa-check"></i> Saved`
}

const alertSaveFail = errorMessage => {
    document.getElementById("save-result").classList.add("fail");
    document.getElementById("save-result").innerHTML = `<i class="fas fa-exclamation-triangle"></i> Failed to Save`;

    document.getElementById("save-error-message").innerHTML = errorMessage;
}