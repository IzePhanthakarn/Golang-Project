const baseURL = 'https://706b-119-42-71-245.ap.ngrok.io/course';
fetch(baseURL).then(function (response){
    return response.json();
}).then(function (data){
    appemdData(data);
}).catch(function (err){
    console.log('error: '+ err);
});

function appemdData(data){
    var mainContainer = document.getElementById("myData");
    for (var i =0;i< data.length; i++){
        var div = document.createElement("div");
        div.innerHTML = 'CourseID: ' + data[i].ID + ' ' + data[i].Name + ' ' + data[i].Price + ' ' + data[i].Instructor;
        mainContainer.appendChild(div);
    }
}