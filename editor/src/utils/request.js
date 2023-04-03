
function Post(url, methon, formData) {
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


function PostGetBlob(url, methon, name) {
  return new Promise(function (resolve, reject) {

    var headers = new Headers();
    headers.append("Content-Type", "application/x-www-form-urlencoded");
    headers.append("FileName", str2utf8(name));

    fetch(url + "/" + methon, {
      method: "POST",
      mode: "cors",
      headers: headers,
    })
      .then((response) => {
        if (response.ok) {
          return response.blob();
        } else {
          reject({ status: response.status });
        }
      })
      .then((response) => {
        resolve({
          blob: response
        });
      })
      .catch((err) => {
        reject({ status: -1 });
      });
  });
}

function fetch_timeout(fecthPromise, timeout = 5000, controller) {
  let abort = null;
  let abortPromise = new Promise((resolve, reject) => {
    abort = () => {
      return reject({
        code: 504,
        message: '请求超时!',
      });
    };
  });
  // 最快出结果的promise 作为结果
  let resultPromise = Promise.race([fecthPromise, abortPromise]);
  setTimeout(() => {
    abort();
    controller.abort();
  }, timeout);

  return resultPromise.then((res) => {
    clearTimeout(timeout);
    return res;
  });
}


function CheckHealth(url) {

  let controller = new AbortController();
  const signal = controller.signal;

  return fetch_timeout(
    fetch(url + "/health", {
      method: "GET",
      mode: 'cors',
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      signal
    }), 3000, controller
  ).then(function (response) {
    if (response.ok) {
      return ({ code: 200 })
    } else {
      return ({ code: response.status });
    }
  }).catch(function (e) {
    console.error(e)
    return ({ code: 404 })
  })

}

var str2utf8 = window.TextEncoder ? function (str) {
  var encoder = new TextEncoder('utf8');
  var bytes = encoder.encode(str);
  var result = '';
  for (var i = 0; i < bytes.length; ++i) {
    result += String.fromCharCode(bytes[i]);
  }
  return result;
} : function (str) {
  return eval('\'' + encodeURI(str).replace(/%/gm, '\\x') + '\'');
}

function PostBlob(url, methon, name, data) {
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
        console.info(err)
        reject({ status: -1 });
      });
  });
}

module.exports = { Post, PostGetBlob, CheckHealth, PostBlob };