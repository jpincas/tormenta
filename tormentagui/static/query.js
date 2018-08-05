const runQuery = entity => {
    fill(buildQueryURL(entity, document.getElementById("query-form")), `results`);
}

const buildQueryURL = (entity, form) => {
    let url = `/${entity}`;
    let queryKeys = ["from", "to", "limit", "offset", "index", "match", "start", "end"];

    // Reverse (and set up queries with ?)
    let reverse = form.elements["reverse"].checked
    if (reverse) {
        url = `${url}?reverse=true`;
    } else {
        url = `${url}?reverse=false`;
    }

    // Rest of the query keys
    queryKeys.map(
        key => {
            let val = form.elements[key].value;
            if (val) { url = `${url}&${key}=${val}`; }
        }
    );

    return url;
}
