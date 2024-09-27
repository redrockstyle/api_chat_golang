const followModal = document.getElementById('follow-chat-modal');
const followCloseBtn = document.getElementById('follow-close-btn');
const followBtn = document.getElementById('follow-btn');
const followChatBtn = document.getElementById('follow-chat-btn');
                    
// Function to show the follow chat modal
function showFollowModal() {
    followModal.style.display = 'block';
    // document.body.style.overflow = 'hidden'; // Disable scrolling
}
                    
// Function to hide the follow chat modal
function hideFollowModal() {
    followModal.style.display = 'none';
    // document.body.style.overflow = 'auto'; // Enable scrolling
}
                   
// Event listener for the follow chat button
followChatBtn.addEventListener('click', showFollowModal);
// Event listener for the close button in the follow chat modal
followCloseBtn.addEventListener('click', hideFollowModal);

// Event listener for the follow button in the follow chat modal
followBtn.addEventListener('click', function() {
    const chatName = document.getElementById('follow-chat-name-input').value;
    if (chatName.trim() !== '') {
        // Perform the necessary actions to follow the chat (e.g., send a request to the server)
        console.log('Following chat:', chatName);

        fetch('/api/chat/follow/' + chatName, { method: 'POST' })
        .catch((error) => { console.log('Error:', error); });

        hideFollowModal();
    } else {
        alert('Please enter a chat name.');
    }
});




// Get the modal elements
const modal = document.getElementById('create-chat-modal');
const closeBtn = document.getElementById('close-btn');
const createBtn = document.getElementById('create-btn');
const createChatBtn = document.getElementById('create-chat-btn');

// Function to show the modal
function showModal() {
  modal.style.display = 'block';
  // document.body.style.overflow = 'hidden'; // Disable scrolling
}

// Function to hide the modal
function hideModal() {
  modal.style.display = 'none';
  // document.body.style.overflow = 'auto'; // Enable scrolling
}

// Event listener for the create chat button
createChatBtn.addEventListener('click', showModal);

// Event listener for the close button
closeBtn.addEventListener('click', hideModal);

// Event listener for the create button
createBtn.addEventListener('click', function() {
  const chatName = document.getElementById('chat-name-input').value;
  if (chatName.trim() !== '') {
    // Perform the necessary actions to create the chat (e.g., send a request to the server)
    console.log('Creating chat:', chatName);
    fetch('/api/chat/create/' + chatName, { method: 'POST' })
    .then(() => { window.location.href = `/chat/${chatName}`; })
    .catch(error => { console.log('Error creating chat:', error)});

    hideModal();
  } else {
    alert('Please enter a chat name.');
  }
});