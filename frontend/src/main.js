import './style.css';
import './app.css';

import logo from './assets/images/logo-universal.png';
import { CallLLM, CancelLLMRequest } from '../wailsjs/go/main/App';

document.getElementById('logo').src = logo;

let userInputElement = document.getElementById('userInput');
let messagesElement = document.getElementById('messages');
let sendBtnElement = document.getElementById('sendBtn');

userInputElement.focus();

// Setup the send message function
window.sendMessage = function () {
  // Get user input
  let userInput = userInputElement.value.trim();

  // Check if the input is empty
  if (userInput === '') return;

  // Disable input and button while processing
  userInputElement.disabled = true;
  sendBtnElement.disabled = true;
  sendBtnElement.textContent = 'Thinking...';

  // Add user message to chat
  addMessage('user', userInput);

  // Clear input
  userInputElement.value = '';

  // Call the LLM
  try {
    CallLLM(userInput)
      .then((result) => {
        // Add assistant response to chat
        addMessage('assistant', result);
      })
      .catch((err) => {
        console.error(err);
        addMessage(
          'assistant',
          'Sorry, I encountered an error: ' + err.message
        );
      })
      .finally(() => {
        // Re-enable input and button
        userInputElement.disabled = false;
        sendBtnElement.disabled = false;
        sendBtnElement.textContent = 'Send';
        userInputElement.focus();
      });
  } catch (err) {
    console.error(err);
    addMessage('assistant', 'Sorry, I encountered an error: ' + err.message);
    // Re-enable input and button
    userInputElement.disabled = false;
    sendBtnElement.disabled = false;
    sendBtnElement.textContent = 'Send';
    userInputElement.focus();
  }
};

// Setup the cancel message function
window.cancelMessage = function () {
  CancelLLMRequest().then((result) => {
    addMessage('assistant', result);
    userInputElement.disabled = false;
    sendBtnElement.disabled = false;
    sendBtnElement.textContent = 'Send';
    userInputElement.focus();
  });
};

// Function to add a message to the chat
function addMessage(sender, content) {
  const messageDiv = document.createElement('div');
  messageDiv.className = `message ${sender}`;
  messageDiv.textContent = content;
  messagesElement.appendChild(messageDiv);
  messagesElement.scrollTop = messagesElement.scrollHeight;
}

// Allow sending message with Enter key (but not Shift+Enter)
userInputElement.addEventListener('keydown', function (e) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault();
    sendMessage();
  }
});
