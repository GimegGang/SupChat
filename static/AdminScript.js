document.addEventListener('DOMContentLoaded', function() {
    const messageForm = document.getElementById("messageForm");
    const messageInput = document.getElementById("messageInput");
    const messageList = document.getElementById("messageList");
    const but = document.getElementById("but");

    messageForm.addEventListener("submit", event => {
        event.preventDefault();

        const message = messageInput.value;

        fetch(document.URL + "/send", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({message: message})
        }).then(response => {
            if(!response.ok){
                throw new Error('Network response was not ok');
            }
            return response.json();
        }).then(data => {
            console.log('Message sent:', data);

            const newMessageItem = document.createElement("li");
            newMessageItem.innerHTML = `<strong>admin</strong><p>${message}</p>`;
            messageList.appendChild(newMessageItem);

            messageInput.value = '';

            fetchNewMessages();
        }).catch((error) => {
            console.error('Error:', error);
        });
    });

    document.getElementById('closeTicket').onclick = () => {
        document.location.href = document.URL + "/closeTicket";
    };

    function fetchNewMessages(){
        const ticketId = document.getElementById("but").value;
        if (!ticketId) {
            console.error("Ticket ID is not set");
            return;
        }

        fetch(`/getMessages/${ticketId}`)
            .then(response => {
                if(!response.ok){
                    throw new Error("Server Error");
                }
                return response.json();
            })
            .then(data => {
                if (data.messages && Array.isArray(data.messages)) {
                    updateMessageList(data.messages);
                } else {
                    console.error("Invalid data format:", data);
                }
            })
            .catch(error => {
                console.error("Error getting messages:", error);
            });
    }

    function updateMessageList(messages) {
        messageList.innerHTML = '';
        messages.forEach(msg => {
            const mesItem = document.createElement("li");
            mesItem.innerHTML = `<strong>${msg.From}</strong><p>${msg.Message}</p>`;
            messageList.appendChild(mesItem);
        });
    }
    if (but.value) {
        setInterval(fetchNewMessages, 1000);
        fetchNewMessages();
    }
});