export function Post(url, methon, formData) {
  return new Promise(function (resolve, reject) {
    fetch(url + "/" + methon, {
      method: "POST",
      mode: "cors",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: JSON.stringify(formData),
    })
      .then((response) => {
        if (response.ok) {
          return response.json();
        } else {
          reject({ status: response.status });
        }
      })
      .then((response) => {
        resolve(response);
      })
      .catch((err) => {
        reject({ status: -1 });
      });
  });
}

var str2utf8 = window.TextEncoder ? function(str) {
  var encoder = new TextEncoder('utf8');
  var bytes = encoder.encode(str);
  var result = '';
  for(var i = 0; i < bytes.length; ++i) {
      result += String.fromCharCode(bytes[i]);
  }
  return result;
} : function(str) {
  return eval('\''+encodeURI(str).replace(/%/gm, '\\x')+'\'');
}

export function PostBlob(url, methon, name, data) {
  return new Promise(function (resolve, reject) {
    var headers = new Headers();
    headers.append("Content-Type", "text/html");
    headers.append("FileName", str2utf8(name));

    fetch(url + "/" + methon, {
      method: "POST",
      mode: "cors",
      headers: headers,
      body: data,
    })
      .then((response) => {
        if (response.ok) {
          return response.json();
        } else {
          reject({ status: response.status });
        }
      })
      .then((response) => {
        resolve(response);
      })
      .catch((err) => {
        reject({ status: -1 });
      });
  });
}
