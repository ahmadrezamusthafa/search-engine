function openModal() {
    resetIndexForm();
    document.getElementById("indexModal").style.display = "block";
}

function closeModal() {
    document.getElementById("indexModal").style.display = "none";
}

function generateDynamicId() {
    const timestamp = new Date().getTime();
    const randomValue = Math.floor(Math.random() * 1000);
    return `id:${timestamp}:${randomValue}`;
}

function validateIndexForm() {
    let isValid = true;

    document.querySelectorAll('.error-message').forEach(el => el.style.display = 'none');
    document.querySelectorAll('input, textarea').forEach(el => el.classList.remove('invalid'));

    const contentString = document.getElementById('contentString').value.trim();
    const contentObject = document.getElementById('contentObject').value.trim();
    const objectIndexes = document.getElementById('objectIndexes').value.trim();

    if (!contentString && !contentObject) {
        document.getElementById('contentStringError').textContent = "Content String or Content Object is required. At least one.";
        document.getElementById('contentStringError').style.display = 'block';
        document.getElementById('contentString').classList.add('invalid');

        document.getElementById('contentObjectError').textContent = "Content String or Content Object is required. At least one.";
        document.getElementById('contentObjectError').style.display = 'block';
        document.getElementById('contentObject').classList.add('invalid');
        isValid = false;
    }

    if (contentObject) {
        try {
            JSON.parse(contentObject);
        } catch {
            document.getElementById('contentObjectError').textContent = "Content Object must be a valid JSON.";
            document.getElementById('contentObjectError').style.display = 'block';
            document.getElementById('contentObject').classList.add('invalid');
            isValid = false;
        }

        if (!objectIndexes) {
            document.getElementById('objectIndexesError').textContent = "Indexes are required.";
            document.getElementById('objectIndexesError').style.display = 'block';
            document.getElementById('objectIndexes').classList.add('invalid');
            isValid = false;
        }
    }

    return isValid;
}

function resetIndexForm() {
    document.getElementById('contentString').value = '';
    document.getElementById('contentObject').value = '';
    document.getElementById('objectIndexes').value = '';
    document.getElementById('stopWords').value = '';

    document.getElementById('contentString').classList.remove('error');
    document.getElementById('contentObject').classList.remove('error');
    document.getElementById('objectIndexes').classList.remove('error');
    document.getElementById('stopWords').classList.remove('error');
}

function submitIndex() {
    if (!validateIndexForm()) {
        return;
    }

    const contentString = document.getElementById('contentString').value;
    let contentObject = document.getElementById('contentObject').value;
    const objectIndexes = document.getElementById('objectIndexes').value.split(',');
    const stopWords = document.getElementById('stopWords').value.split(',');

    let dynamicId = generateDynamicId();
    if (!contentObject) {
        contentObject = null;
    }

    const data = {
        id: dynamicId, content: {
            string: contentString,
            object: JSON.parse(contentObject),
            object_indexes: objectIndexes.map(item => item.trim())
        }, stop_words: stopWords.map(word => word.trim())
    };

    fetch('/index', {
        method: 'POST', headers: {
            'Content-Type': 'application/json'
        }, body: JSON.stringify(data)
    })
        .then(response => response.json())
        .then(data => {
            document.getElementById('resultOutput').textContent = JSON.stringify(data, null, 2);
            closeModal();
        })
        .catch(error => {
            document.getElementById('resultOutput').textContent = 'Error: ' + error;
            closeModal();
        });
}

function submitSearch() {
    const queryStrings = document.getElementById('queryStrings').value.split(',');
    const queryParams = queryStrings.map(query => `query=${query.trim()}`).join('&');

    fetch(`/search?${queryParams}`, {
        method: 'GET'
    })
        .then(response => response.json())
        .then(data => {
            document.getElementById('resultOutput').textContent = JSON.stringify(data, null, 2);
        })
        .catch(error => {
            document.getElementById('resultOutput').textContent = 'Error: ' + error;
        });
}

