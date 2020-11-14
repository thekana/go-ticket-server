function currentTime() {
	var date = new Date(); /* creating object of Date class */
	var hour = date.getHours();
	var min = date.getMinutes();
	var sec = date.getSeconds();
	hour = updateTime(hour);
	min = updateTime(min);
	sec = updateTime(sec);
	document.getElementById("clock").innerText = hour + " : " + min + " : " + sec; /* adding time to the div */
	  var t = setTimeout(function(){ currentTime() }, 1000); /* setting timer */
  }
  
  function updateTime(k) {
	if (k < 10) {
	  return "0" + k;
	}
	else {
	  return k;
	}
  }
  
  currentTime(); /* calling currentTime() function to initiate the process */