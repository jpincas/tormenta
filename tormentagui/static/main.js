
const fill = (url, target, callback) => {
    fetch(url)
        .then(function (response) {
            return response.text();
        })
        .then(function (contents) {
            morphdom(document.getElementById(target), `<div id="${target}">${contents}</div>`);
            if (callback) {
                callback();
            }
        });
}



const saveEntity = (entity, id) => {
    const data = document.getElementById("entity-json").innerHTML;
    postData(`/api/${entity}/${id}`, data)
        .then(response => saveResult(response, entity));
}

// Ajax methods

const postData = (url = ``, data = "") => {
    return fetch(url, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json; charset=utf-8",

        },
        body: data,
    })
};

// Save Result

// saveResult executes the necessary actions after an entity save
const saveResult = (response, entity) => {
    // For an OK response
    if (response.ok) {
        // Decode the JSON
        response.json().then(
            saved => {
                // Using the returned ID,
                // get the new/edited record
                // and refresh the detail view
                fill(`/${entity}/${saved.ID}`, 'detail-view', alertSaveSuccess)

                //And finally refresh the list view
                fill(`/${entity}`, 'list-view')
            }
        )
    } else {
        response.json().then(
            error => alertSaveFail(error.errorMessage)
        )
    }
}

// Alerts

const alertSaveSuccess = () => {
    document.getElementById("save-result").classList.add("success");
    document.getElementById("save-result").innerHTML = `<i class="fas fa-check"></i> Saved`


}

const alertSaveFail = errorMessage => {
    document.getElementById("save-result").classList.add("fail");
    document.getElementById("save-result").innerHTML = `<i class="fas fa-exclamation-triangle"></i> Failed to Save! `;

    document.getElementById("save-error-message").innerHTML = errorMessage;
}