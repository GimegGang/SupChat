document.getElementById("messageForm").addEventListener("submit", (event) => {
    event.preventDefault()

    const messageInput = document.getElementById("messageInput")
    const message = messageInput.value
    const but = document.getElementById("but")
    let tId = ""

    if(but.value === ""){
        tId = "new"
    }else{
        tId = but.value
    }

    fetch("/sendUserMessage", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({message: message, tId: tId})
    }).then(response => {
        if(!response.ok){
            throw new Error('Network response was not ok')
        }
        return response.json()
    }).then(data => {
        console.log('Message sent:', data);

        const messageList = document.getElementById("messageList");
        const newMessageItem = document.createElement("li");

        // Устанавливаем содержимое нового элемента списка
        newMessageItem.innerHTML = `<strong>User</strong><p>${message}</p>`;

        // Добавляем новый элемент в список
        messageList.appendChild(newMessageItem);

        messageInput.value = '';
        but.value = data.ticketId;
    })
        .catch((error) => {
            console.error('Error:', error);
        });
})

function fetchNewMessages(){
    const ticketId = document.getElementById("but").value;
    fetch(`/getMessages/${ticketId}`).then(response => {
        if(!response.ok){
            throw new Error("Server Error")
        }
        return response.json()
    }).then(data => {
        if(!data.closed){
            const messageList = document.getElementById("messageList");
            messageList.innerHTML = '';
            data.messages.forEach(msg => {
                const mesItem = document.createElement("li")
                mesItem.innerHTML = `<strong>${msg.From}</strong><p>${msg.Message}</p>`
                messageList.appendChild(mesItem)
            })
        }else{
            if(document.getElementById("messageForm")){
                document.getElementById("messageForm").remove()
                document.getElementById("messageList").innerHTML = `<h2>Ticket is closed by Administrator</h2>`
            }
        }
    }).catch(error => {
        console.error("error get mes", error)
    })
}

if (document.getElementById("but").value !== "new") {
    setInterval(fetchNewMessages, 1000);
}