export function Post(url, methon, formData) {
  return new Promise(function (resolve, reject) {
    fetch(url + methon, {
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

export function PostBlob(url, methon, name, data) {
  return new Promise(function (resolve, reject) {
    fetch(url + methon, {
      method: "POST",
      mode: "cors",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
        FileName: name,
      },
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
