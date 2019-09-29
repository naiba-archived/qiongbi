$(function () {

  // 一言
  var xhr = new XMLHttpRequest();
  xhr.open('get', 'https://v1.hitokoto.cn');
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4) {
      var data = JSON.parse(xhr.responseText);
      var hitokoto = document.getElementById('note');
      hitokoto.setAttribute('placeholder', data.hitokoto);
    }
  }
  xhr.send();

  $("#email").change(function () {
    var reg = /^([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+\.[a-zA-Z]{2,20}$/;
    var email = $(this).val()
    if (!email == "") {
      if (!reg.test(email)) {
        showmsg("请输入有效的Email地址", "error")
        $(this).val("").focus()
        return false;
      }
    }
  })

  $(".amount").click(function () {
    $(this).addClass("checked").siblings().removeClass("checked")
    $("#customAmount").val(1)
    $("#amount").val($(this).data("amount"))
  })

  $("#customAmount").keyup(function () {
    var amount = $(this).val()
    $(this).val(amount)
    $("#amount").val(amount)
  }).bind("paste", function () { $(this).val($(this).val()) })

})

function checkSubmit() {
  let flag = false
  if (!$("input#name").val()) {
    showmsg("做好事请留名", 'error')
    return flag
  }

  if (!$("input#email").val()) {
    showmsg("请输入正确的邮箱", 'error')
    return flag
  }

  if ($("input#amount").val() < 1) {
    showmsg("请多捐点吧", 'error')
    return flag
  }

  flag = true
  return flag
}

function showmsg(msg, type) {
  msgstr = '<div class="msg-wrap"><div class="msg"><span class="ico">'
  if (type == "error") {
    msgstr = msgstr + "❎"
  } else {
    msgstr = msgstr + "✅"
  }
  msgstr = msgstr + '</span><span class="txt">' + msg + '</span></div></div>'
  $(".alms-form").append(msgstr)
  setTimeout(function () {
    $(".msg-wrap").remove()
  }, 1500)
}