function makeUL(array) {
	var list = document.createElement('ul');

	array = array.reverse();

	for(i = 0; i < array.length; i++) {
		var item = document.createElement('li');
		a = document.createElement('a');
		a.href = array[i];
		a.innerHTML = decodeURIComponent(array[i].replace(/\+/g, " "));
		item.appendChild(a);
		list.appendChild(item);
	}

	return list;
};

function giphy(category) {
	var url = 'http://api.giphy.com';
	var query = '/v1/gifs/random';
	var params = '?limit=1&api_key=dc6zaTOxFJmzC';

	switch(category) {
	case "0":
		params += '&tag=cats';
		break;
	case "1":
		params += '&tag=chuck+norris';
		break;
	case "3":
		params += '&tag=numbers';
		break;
	case "4":
		params += '&tag=archer';
		break;
	};

	var xmlHttp = null;
	xmlHttp = new XMLHttpRequest();
	xmlHttp.open("GET", url + query + params, false);
	xmlHttp.send(null);

	if (xmlHttp.status == 200) {
		var gif = document.getElementById('gif');
		var urls = JSON.parse(xmlHttp.responseText);
		if (urls == null) {
			return false;
		}
		gif.innerHTML = '<img src="'+urls.data.fixed_height_downsampled_url+'"/>';
	}
};

function lengthen() {
	var url = document.getElementById('url').value;
	var category = document.getElementById('category').value;

	if (!url.match('https?://')) {
		url = 'http://' + url;
	}

	var xmlHttp = null;
	xmlHttp = new XMLHttpRequest();
	xmlHttp.open("GET", '/lengthen?url=' + encodeURIComponent(url) + '&category=' + category, false);
	xmlHttp.send(null);

	var cruft = document.getElementById('cruft');
	if (xmlHttp.status == 200) {
		cruft.innerHTML = '<a href="'+xmlHttp.responseText+'">'+decodeURIComponent(xmlHttp.responseText.replace(/\+/g, " "))+'</a>';
	} else {
		cruft.innerHTML = xmlHttp.responseText;
	}

	document.getElementById('cruft').style.display = 'block';
	latest();
	giphy(category);
	return false;
};

function latest() {
	var xmlHttp = null;
	xmlHttp = new XMLHttpRequest();
	xmlHttp.open("GET", '/latest', false);
	xmlHttp.send(null);

	if (xmlHttp.status == 200) {
		var urls = JSON.parse(xmlHttp.responseText);
		if (urls == null) {
			return false;
		}

		var list = document.getElementById('list');
		while (list.lastChild) {
			list.removeChild(list.lastChild);
		}
		list.appendChild(makeUL(urls));     
		list.style.display = 'block';
	}

	return false;
};

