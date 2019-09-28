$(function () {
  
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
    $("#customAmount").val(0)
    $("#amount").val($(this).data("amount"))
  })

  $("#customAmount").keyup(function () {
    var amount = $(this).val()
    $(this).val(amount)
    $("#amount").val(amount)
  }).bind("paste", function () { $(this).val($(this).val()) })

})

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