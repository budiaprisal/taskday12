function submitData() {
  let name = document.getElementById("input-name").value;
  let email = document.getElementById("input-email").value;
  let phone = document.getElementById("input-phone").value;
  let subject = document.getElementById("input-subject").value;
  let message = document.getElementById("input-message").value;

  if (email == "") {
    return alert("name wajib diisi");
  } else if (email == "") {
    return alert("email wajib diisi");
  } else if (phone == "") {
    return alert("phone wajib diisi");
  } else if (subject == "") {
    return alert("subject wajib diisi");
  } else if (message == "") {
    return alert("message wajib diisi");
  }

  console.log(name);
  console.log(email);
  console.log(phone);
  console.log(subject);
  console.log(message);

  let emailReceiver = "budiaprisal@mail.com";

  let a = document.createElement("a");
  a.href = `mailto:${emailReceiver}?subject=${subject}&body= Hai Perkenalkan  nama saya ${name} , ${message}, silahkan kontak saya dengan email ${email}, telp ${phone}`;
  a.click();
}
