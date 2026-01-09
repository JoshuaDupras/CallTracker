// Import Wails runtime
let currentUser = null;
let picklists = {};
let modalCallback = null;

// Modal helper functions
function showModal(title, bodyHTML, confirmText = 'OK', confirmCallback = null) {
    document.getElementById('modal-title').textContent = title;
    document.getElementById('modal-body').innerHTML = bodyHTML;
    document.getElementById('modal-confirm-btn').textContent = confirmText;
    document.getElementById('modal-overlay').style.display = 'flex';
    modalCallback = confirmCallback;
    
    // Focus first input if exists and add Enter key handlers
    setTimeout(() => {
        const inputs = document.querySelectorAll('#modal-body input, #modal-body select');
        if (inputs.length > 0) {
            inputs[0].focus();
            
            // Add Enter key handler to all inputs
            inputs.forEach((input, index) => {
                input.addEventListener('keydown', function(e) {
                    if (e.key === 'Enter') {
                        e.preventDefault();
                        // If not the last input, move to next field
                        if (index < inputs.length - 1) {
                            inputs[index + 1].focus();
                        } else {
                            // On last input, trigger confirm
                            confirmModal();
                        }
                    }
                });
            });
        }
    }, 100);
}

function closeModal() {
    document.getElementById('modal-overlay').style.display = 'none';
    document.getElementById('modal-body').innerHTML = '';
    modalCallback = null;
}

function confirmModal() {
    if (modalCallback) {
        modalCallback();
    }
    closeModal();
}

// Initialize app
window.onload = async function() {
    await loadUsers();
};

// Load users for login dropdown
async function loadUsers() {
    try {
        const users = await window.go.main.App.GetActiveUsers();
        const select = document.getElementById('login-name');
        select.innerHTML = '<option value="">Select your name...</option>';
        
        if (users && users.length > 0) {
            users.forEach(user => {
                const option = document.createElement('option');
                const fullName = `${user.first_name} ${user.last_name}`;
                option.value = fullName;
                option.textContent = fullName;
                select.appendChild(option);
            });
        } else {
            showAdminLogin();
            document.getElementById('admin-link-div').style.display = 'none';
        }
    } catch (error) {
        console.error('Failed to load users:', error);
        showAdminLogin();
    }
}

// Called when user selects their name
function userSelected() {
    const select = document.getElementById('login-name');
    const pinGroup = document.getElementById('pin-group');
    const loginBtn = document.getElementById('login-btn');
    const errorDiv = document.getElementById('login-error');
    
    if (select.value) {
        // User selected - show PIN field and login button
        pinGroup.style.display = 'block';
        loginBtn.style.display = 'block';
        errorDiv.textContent = '';
        // Focus PIN field
        setTimeout(() => document.getElementById('login-pin').focus(), 100);
    } else {
        // No user selected - hide PIN and button
        pinGroup.style.display = 'none';
        loginBtn.style.display = 'none';
        document.getElementById('login-pin').value = '';
    }
}

// Login
async function doLogin() {
    let name;
    const textInput = document.getElementById('login-name-text');
    const selectInput = document.getElementById('login-name');
    
    if (textInput && textInput.value) {
        name = textInput.value;
    } else {
        name = selectInput.value;
    }
    
    const pin = document.getElementById('login-pin').value;
    const errorDiv = document.getElementById('login-error');
    const warningDiv = document.getElementById('admin-warning');
    
    if (!name || !pin) {
        errorDiv.textContent = 'Please enter your name and PIN';
        warningDiv.style.display = 'none';
        return;
    }
    
    try {
        currentUser = await window.go.main.App.Login(name, pin);
        errorDiv.textContent = '';
        warningDiv.style.display = 'none';
        showMainMenu();
    } catch (error) {
        errorDiv.textContent = 'Invalid credentials';
        document.getElementById('login-pin').value = '';
        
        // Show admin contact warning
        try {
            const admins = await window.go.main.App.GetAdminUsers();
            if (admins && admins.length > 0) {
                const adminNames = admins.map(a => `${a.first_name} ${a.last_name}`).join(', ');
                warningDiv.innerHTML = `<strong>User not found or incorrect PIN.</strong><br>Please contact a call log administrator: ${adminNames}`;
                warningDiv.style.display = 'block';
            }
        } catch (adminError) {
            console.error('Failed to load admins:', adminError);
        }
    }
}

// Logout
async function doLogout() {
    await window.go.main.App.Logout();
    currentUser = null;
    showScreen('login-screen');
    document.getElementById('login-pin').value = '';
    document.getElementById('pin-group').style.display = 'none';
    document.getElementById('login-btn').style.display = 'none';
    document.getElementById('login-error').textContent = '';
    document.getElementById('admin-warning').style.display = 'none';
    
    // Reset to dropdown if admin text input was showing
    const textInput = document.getElementById('login-name-text');
    if (textInput) {
        textInput.remove();
        document.getElementById('login-name').style.display = 'block';
        document.getElementById('admin-link-div').style.display = 'block';
    }
    document.getElementById('login-name').value = '';
}

// Change current user's PIN
function showChangePIN() {
    if (!currentUser) {
        alert('Not logged in');
        return;
    }
    
    // Hardcoded admin cannot change PIN
    if (currentUser.id === 0) {
        alert('Cannot change hardcoded admin PIN');
        return;
    }
    
    const modalBody = `
        <div class="form-group">
            <label>Current PIN</label>
            <input type="password" id="change-pin-old" maxlength="4" class="form-control">
        </div>
        <div class="form-group">
            <label>New PIN (4 digits)</label>
            <input type="password" id="change-pin-new" maxlength="4" class="form-control">
        </div>
        <div class="form-group">
            <label>Confirm New PIN</label>
            <input type="password" id="change-pin-confirm" maxlength="4" class="form-control">
        </div>
    `;
    
    showModal('Change Your PIN', modalBody, 'Change PIN', () => {
        const oldPIN = document.getElementById('change-pin-old').value;
        const newPIN = document.getElementById('change-pin-new').value;
        const confirmPIN = document.getElementById('change-pin-confirm').value;
        
        if (!oldPIN) {
            alert('Please enter your current PIN');
            return;
        }
        
        if (!newPIN || newPIN.length !== 4) {
            alert('New PIN must be 4 digits');
            return;
        }
        
        if (newPIN !== confirmPIN) {
            alert('New PINs do not match');
            return;
        }
        
        window.go.main.App.ChangePIN(oldPIN, newPIN)
            .then(() => {
                alert('PIN changed successfully!');
            })
            .catch(error => {
                alert('Failed to change PIN: ' + error);
            });
    });
}

// Show admin login - allow manual entry
function showAdminLogin() {
    const select = document.getElementById('login-name');
    const pinGroup = document.getElementById('pin-group');
    const loginBtn = document.getElementById('login-btn');
    const adminLinkDiv = document.getElementById('admin-link-div');
    const backToUsersDiv = document.getElementById('back-to-users-div');
    
    // Check if text input already exists
    let textInput = document.getElementById('login-name-text');
    if (textInput) {
        // Already in admin mode, just focus PIN
        document.getElementById('login-pin').focus();
        return;
    }
    
    // Hide the dropdown and admin link
    select.style.display = 'none';
    adminLinkDiv.style.display = 'none';
    backToUsersDiv.style.display = 'block';
    
    // Create text input to replace select
    textInput = document.createElement('input');
    textInput.type = 'text';
    textInput.id = 'login-name-text';
    textInput.value = 'Admin User';
    textInput.placeholder = 'Enter full name';
    textInput.className = select.className;
    select.parentNode.insertBefore(textInput, select);
    
    // Show PIN field and login button
    pinGroup.style.display = 'block';
    loginBtn.style.display = 'block';
    
    // Focus the PIN field
    document.getElementById('login-pin').focus();
}

// Back to regular user login
function backToUserLogin() {
    const select = document.getElementById('login-name');
    const textInput = document.getElementById('login-name-text');
    const pinGroup = document.getElementById('pin-group');
    const loginBtn = document.getElementById('login-btn');
    const adminLinkDiv = document.getElementById('admin-link-div');
    const backToUsersDiv = document.getElementById('back-to-users-div');
    const errorDiv = document.getElementById('login-error');
    const warningDiv = document.getElementById('admin-warning');
    
    // Remove text input if it exists
    if (textInput) {
        textInput.remove();
    }
    
    // Show dropdown and admin link, hide back link
    select.style.display = 'block';
    select.value = '';
    adminLinkDiv.style.display = 'block';
    backToUsersDiv.style.display = 'none';
    
    // Hide PIN field and login button
    pinGroup.style.display = 'none';
    loginBtn.style.display = 'none';
    document.getElementById('login-pin').value = '';
    
    // Clear errors
    errorDiv.textContent = '';
    warningDiv.style.display = 'none';
}

// Show main menu
function showMainMenu() {
    if (!currentUser) {
        console.error('No current user set');
        return;
    }
    
    const fullName = `${currentUser.first_name || ''} ${currentUser.last_name || ''}`;
    const isAdmin = currentUser.is_admin || false;
    
    document.getElementById('current-user').textContent = `Logged in as: ${fullName.trim()}`;
    
    // Show/hide admin section
    if (isAdmin) {
        document.getElementById('admin-section').style.display = 'block';
    } else {
        document.getElementById('admin-section').style.display = 'none';
    }
    
    showScreen('menu-screen');
}

// Navigation functions
function showScreen(screenId) {
    document.querySelectorAll('.screen').forEach(s => s.style.display = 'none');
    document.getElementById(screenId).style.display = 'block';
}

function backToMenu() {
    showMainMenu();
}

// New Call
let currentWizardStep = 1;
const totalWizardSteps = 12;

async function showNewCall() {
    showScreen('newcall-screen');
    currentWizardStep = 1;
    await loadPicklists();
    await loadResponders();
    clearNewCallForm();
    
    // Set default date to today and time to current time (after clearing form)
    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0');
    const day = String(now.getDate()).padStart(2, '0');
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    
    document.getElementById('q-dispatched-date').value = `${year}-${month}-${day}`;
    document.getElementById('q-dispatched-time').value = `${hours}:${minutes}`;
    
    // Generate incident number immediately with default date
    updateDispatchedValue();
    
    updateWizardDisplay();
}

async function loadResponders() {
    try {
        const users = await window.go.main.App.GetActiveUsers();
        const respondersDiv = document.getElementById('responders-checkboxes');
        respondersDiv.innerHTML = '';
        
        users.forEach(user => {
            const label = document.createElement('label');
            label.style.display = 'flex';
            label.style.alignItems = 'center';
            label.style.cursor = 'pointer';
            const fullName = `${user.first_name} ${user.last_name}`;
            label.innerHTML = `
                <input type="checkbox" name="responders" value="${user.id}" style="margin-right: 10px; width: 20px; height: 20px;">
                <span style="font-size: 1.1em;">${fullName}</span>
            `;
            respondersDiv.appendChild(label);
        });
    } catch (error) {
        console.error('Failed to load responders:', error);
    }
}

function setDispatchedDate(option) {
    const now = new Date();
    if (option === 'yesterday') {
        now.setDate(now.getDate() - 1);
    }
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0');
    const day = String(now.getDate()).padStart(2, '0');
    const dateString = `${year}-${month}-${day}`;
    
    const dateInput = document.getElementById('q-dispatched-date');
    dateInput.value = dateString;
    
    // If time is not set, set to current time
    const timeInput = document.getElementById('q-dispatched-time');
    if (!timeInput.value) {
        const hours = String(new Date().getHours()).padStart(2, '0');
        const minutes = String(new Date().getMinutes()).padStart(2, '0');
        timeInput.value = `${hours}:${minutes}`;
    }
    
    // Update combined value and generate incident number
    updateDispatchedValue();
}

async function loadPicklists() {
    const categories = ['call_type', 'mutual_aid', 'mutual_aid_agencies', 'town', 'apparatus'];
    
    for (const category of categories) {
        try {
            const items = await window.go.main.App.GetPicklistByCategory(category);
            picklists[category] = items;
            
            // Handle apparatus checkboxes
            if (category === 'apparatus') {
                const apparatusDiv = document.getElementById('apparatus-checkboxes');
                if (apparatusDiv) {
                    apparatusDiv.innerHTML = '';
                    items.forEach(item => {
                        const label = document.createElement('label');
                        label.style.display = 'flex';
                        label.style.alignItems = 'center';
                        label.style.cursor = 'pointer';
                        label.innerHTML = `
                            <input type="checkbox" name="apparatus" value="${item.id}" style="margin-right: 10px; width: 20px; height: 20px;">
                            <span style="font-size: 1.1em;">${item.value}</span>
                        `;
                        apparatusDiv.appendChild(label);
                    });
                }
                continue;
            }
            
            // Handle mutual aid agencies with datalist
            if (category === 'mutual_aid_agencies') {
                const input = document.getElementById('q-mutual-aid-agencies-input');
                const datalistId = 'q-mutual-aid-agencies-list';
                
                if (input) {
                    let datalist = document.getElementById(datalistId);
                    if (!datalist) {
                        datalist = document.createElement('datalist');
                        datalist.id = datalistId;
                        input.parentNode.appendChild(datalist);
                        input.setAttribute('list', datalistId);
                    }
                    
                    datalist.innerHTML = '';
                    items.forEach(item => {
                        const option = document.createElement('option');
                        option.value = item.value;
                        datalist.appendChild(option);
                    });
                }
                continue;
            }
            
            const fieldId = 'q-' + category.replace('_', '-');
            const input = document.getElementById(fieldId);
            const datalistId = fieldId + '-list';
            
            if (input) {
                // Set up datalist
                let datalist = document.getElementById(datalistId);
                if (!datalist) {
                    datalist = document.createElement('datalist');
                    datalist.id = datalistId;
                    input.parentNode.appendChild(datalist);
                    input.setAttribute('list', datalistId);
                }
                
                datalist.innerHTML = '';
                items.forEach(item => {
                    const option = document.createElement('option');
                    option.value = item.value;
                    datalist.appendChild(option);
                });
            }
        } catch (error) {
            console.error(`Failed to load ${category}:`, error);
        }
    }
}

function updateDispatchedValue() {
    const dateValue = document.getElementById('q-dispatched-date').value;
    const timeValue = document.getElementById('q-dispatched-time').value;
    
    if (dateValue && timeValue) {
        const combined = `${dateValue}T${timeValue}`;
        document.getElementById('dispatched').value = combined;
        updateSummary('dispatched', combined);
        
        // Generate incident number from date
        const year = new Date(dateValue).getFullYear();
        loadNextCallNumber(year);
    }
}

function handleMutualAidChange() {
    const mutualAidValue = document.getElementById('q-mutual-aid').value;
    const agenciesCard = document.querySelector('[data-step="3.5"]');
    
    if (mutualAidValue === 'Yes') {
        agenciesCard.style.display = 'block';
    } else {
        agenciesCard.style.display = 'none';
        // Clear agencies
        document.getElementById('selected-agencies').innerHTML = '';
        document.getElementById('mutual-aid-agencies').value = '';
    }
}

let selectedAgencies = [];

function addMutualAidAgency() {
    const input = document.getElementById('q-mutual-aid-agencies-input');
    const agencyName = input.value.trim();
    
    if (!agencyName) {
        alert('Please enter an agency name');
        return;
    }
    
    // Check if already added
    if (selectedAgencies.includes(agencyName)) {
        alert('This agency has already been added');
        return;
    }
    
    // Add to array
    selectedAgencies.push(agencyName);
    
    // Update display
    updateSelectedAgenciesDisplay();
    
    // Clear input
    input.value = '';
}

function removeMutualAidAgency(agencyName) {
    selectedAgencies = selectedAgencies.filter(a => a !== agencyName);
    updateSelectedAgenciesDisplay();
}

function updateSelectedAgenciesDisplay() {
    const container = document.getElementById('selected-agencies');
    
    if (selectedAgencies.length === 0) {
        container.innerHTML = '<p style="color: #999; font-style: italic;">No agencies selected</p>';
        document.getElementById('mutual-aid-agencies').value = '';
        return;
    }
    
    container.innerHTML = selectedAgencies.map(agency => `
        <div class="agency-tag" style="display: inline-block; background: #e3f2fd; padding: 8px 12px; margin: 5px; border-radius: 4px; font-size: 14px;">
            ${agency}
            <button type="button" onclick="removeMutualAidAgency('${agency.replace(/'/g, "\\'")}')" style="margin-left: 8px; background: none; border: none; color: #d32f2f; cursor: pointer; font-weight: bold;">×</button>
        </div>
    `).join('');
    
    // Update hidden input with comma-separated list
    document.getElementById('mutual-aid-agencies').value = selectedAgencies.join(', ');
}

async function loadNextCallNumber(year) {
    try {
        if (!year) {
            year = new Date().getFullYear();
        }
        const nextNumber = await window.go.main.App.GetNextCallNumber(year);
        document.getElementById('incident-number').value = nextNumber;
        document.getElementById('incident-display').textContent = nextNumber;
        document.getElementById('incident-number-display').textContent = nextNumber;
        document.getElementById('incident-number-header').style.display = 'block';
        updateSummary('incident-number', nextNumber);
    } catch (error) {
        console.error('Failed to load next call number:', error);
    }
}

function updateWizardDisplay() {
    // Update progress
    document.getElementById('question-progress').textContent = `Question ${currentWizardStep} of ${totalWizardSteps}`;
    
    // Hide all questions
    document.querySelectorAll('.question-card').forEach(card => {
        card.classList.remove('active');
        // Hide conditional cards that use inline styles
        if (card.getAttribute('data-step') === '3.5') {
            // Always hide step 3.5 first, then show it only if we're on that step
            card.style.display = 'none';
        }
    });
    
    // Show current question
    const currentCard = document.querySelector(`[data-step="${currentWizardStep}"]`);
    if (currentCard) {
        currentCard.classList.add('active');
        
        // If this is the conditional agencies step, make sure it's visible
        if (currentCard.getAttribute('data-step') === '3.5') {
            currentCard.style.display = 'block';
        }
        
        // Default time fields based on previous times
        const field = currentCard.getAttribute('data-field');
        const dispatchedValue = document.getElementById('dispatched').value;
        
        if (dispatchedValue) {
            const dispatchedDate = new Date(dispatchedValue);
            const dateStr = dispatchedDate.toISOString().split('T')[0];
            const timeStr = dispatchedDate.toTimeString().substring(0, 5);
            
            // Default enroute to dispatch date/time
            if (field === 'enroute') {
                const dateInput = document.getElementById('q-enroute-date');
                const timeInput = document.getElementById('q-enroute-time');
                if (!dateInput.value) dateInput.value = dateStr;
                if (!timeInput.value) timeInput.value = timeStr;
            }
            
            // Default on-scene to dispatch date, and time from enroute if available
            if (field === 'on-scene') {
                const dateInput = document.getElementById('q-on-scene-date');
                const timeInput = document.getElementById('q-on-scene-time');
                if (!dateInput.value) dateInput.value = dateStr;
                if (!timeInput.value) {
                    const enrouteValue = document.getElementById('enroute').value;
                    if (enrouteValue) {
                        const enrouteDate = new Date(enrouteValue);
                        timeInput.value = enrouteDate.toTimeString().substring(0, 5);
                    }
                }
            }
            
            // Default clear to dispatch date, and time from on-scene if available
            if (field === 'clear') {
                const dateInput = document.getElementById('q-clear-date');
                const timeInput = document.getElementById('q-clear-time');
                if (!dateInput.value) dateInput.value = dateStr;
                if (!timeInput.value) {
                    const onSceneValue = document.getElementById('on-scene').value;
                    if (onSceneValue) {
                        const onSceneDate = new Date(onSceneValue);
                        timeInput.value = onSceneDate.toTimeString().substring(0, 5);
                    }
                }
            }
        }
        
        // Focus the input
        const input = currentCard.querySelector('input, select, textarea');
        if (input && !input.readOnly) {
            setTimeout(() => input.focus(), 100);
        }
    }
    
    // Update buttons
    const btnPrevious = document.getElementById('btn-previous');
    const btnNext = document.getElementById('btn-next');
    const btnSave = document.getElementById('btn-save');
    
    btnPrevious.style.display = currentWizardStep === 1 ? 'none' : 'inline-block';
    btnNext.style.display = currentWizardStep === totalWizardSteps ? 'none' : 'inline-block';
    btnSave.style.display = currentWizardStep === totalWizardSteps ? 'inline-block' : 'none';
}

function wizardNext() {
    // Get current question's field
    const currentCard = document.querySelector(`[data-step="${currentWizardStep}"]`);
    const field = currentCard.getAttribute('data-field');
    
    // Handle dispatched field specially (has two inputs)
    if (field === 'dispatched') {
        const dateInput = document.getElementById('q-dispatched-date');
        const timeInput = document.getElementById('q-dispatched-time');
        const dateValue = dateInput.value;
        const timeValue = timeInput.value;
        
        // Validate both date and time are filled
        if (!dateValue || !timeValue) {
            alert('Please enter both date and time');
            return;
        }
        
        // Combine and save
        const combined = `${dateValue}T${timeValue}`;
        const dispatchedDateTime = new Date(combined);
        
        // Check if time is in the future
        if (dispatchedDateTime > new Date()) {
            alert('Dispatched time cannot be in the future');
            return;
        }
        
        document.getElementById('dispatched').value = combined;
        updateSummary('dispatched', combined);
        
        // Generate incident number
        const year = new Date(dateValue).getFullYear();
        loadNextCallNumber(year);
    } else if (field === 'mutual-aid') {
        // Handle mutual aid - check the value and show/hide agencies question
        const qInput = document.getElementById('q-' + field);
        const hiddenInput = document.getElementById(field);
        hiddenInput.value = qInput.value;
        updateSummary(field, qInput.value);
        
        // Update agencies card visibility
        handleMutualAidChange();
    } else if (field === 'mutual-aid-agencies') {
        // Agencies are already saved in the hidden input via updateSelectedAgenciesDisplay
        // No additional action needed here
    } else if (field === 'apparatus') {
        // Handle apparatus checkboxes
        const checkedBoxes = document.querySelectorAll('input[name="apparatus"]:checked');
        const apparatusNames = Array.from(checkedBoxes).map(cb => {
            return cb.parentElement.querySelector('span').textContent;
        });
        updateSummary(field, apparatusNames.join(', ') || 'None');
    } else if (field === 'responders') {
        // Handle responders checkboxes
        const checkedBoxes = document.querySelectorAll('input[name="responders"]:checked');
        const responderNames = Array.from(checkedBoxes).map(cb => {
            return cb.parentElement.querySelector('span').textContent;
        });
        updateSummary(field, responderNames.join(', ') || 'None');
    } else if (field === 'enroute' || field === 'on-scene' || field === 'clear') {
        // Handle time fields with separate date and time
        const dateInput = document.getElementById(`q-${field}-date`);
        const timeInput = document.getElementById(`q-${field}-time`);
        const dateValue = dateInput.value;
        const timeValue = timeInput.value;
        
        // If time is provided, date must also be provided
        if (timeValue && !dateValue) {
            alert('Please enter both date and time');
            return;
        }
        
        // If both are provided, combine them
        if (dateValue && timeValue) {
            const combined = `${dateValue}T${timeValue}`;
            const fieldDateTime = new Date(combined);
            
            // Check if time is in the future
            if (fieldDateTime > new Date()) {
                const fieldName = field === 'enroute' ? 'Enroute' : field === 'on-scene' ? 'On scene' : 'Clear';
                alert(`${fieldName} time cannot be in the future`);
                return;
            }
            
            document.getElementById(field).value = combined;
            updateSummary(field, combined);
        } else {
            // Clear if neither provided
            document.getElementById(field).value = '';
            updateSummary(field, '');
        }
    } else {
        // Normal field handling
        const qInput = document.getElementById('q-' + field);
        
        // Save value to hidden input
        const hiddenInput = document.getElementById(field);
        hiddenInput.value = qInput.value;
        
        // Update summary
        updateSummary(field, qInput.value);
        
        // Validate required fields
        const subtitle = currentCard.querySelector('.question-subtitle');
        if (subtitle && subtitle.textContent.includes('Required') && !qInput.value) {
            alert('This field is required');
            return;
        }
    }
    
    // Validate time order
    if (field === 'enroute' || field === 'on-scene' || field === 'clear') {
        const dispatched = document.getElementById('dispatched').value;
        const enroute = document.getElementById('enroute').value;
        const onScene = document.getElementById('on-scene').value;
        const clear = document.getElementById('clear').value;
        
        if (field === 'enroute' && enroute && dispatched) {
            const dispatchedDate = new Date(dispatched);
            const enrouteDate = new Date(enroute);
            if (enrouteDate < dispatchedDate) {
                alert('Enroute time cannot be before dispatched time');
                return;
            }
        }
        
        if (field === 'on-scene' && onScene) {
            if (dispatched) {
                const dispatchedDate = new Date(dispatched);
                const onSceneDate = new Date(onScene);
                if (onSceneDate < dispatchedDate) {
                    alert('On scene time cannot be before dispatched time');
                    return;
                }
            }
            if (enroute) {
                const enrouteDate = new Date(enroute);
                const onSceneDate = new Date(onScene);
                if (onSceneDate < enrouteDate) {
                    alert('On scene time cannot be before enroute time');
                    return;
                }
            }
        }
        
        if (field === 'clear' && clear) {
            if (dispatched) {
                const dispatchedDate = new Date(dispatched);
                const clearDate = new Date(clear);
                if (clearDate < dispatchedDate) {
                    alert('Clear time cannot be before dispatched time');
                    return;
                }
            }
            if (enroute) {
                const enrouteDate = new Date(enroute);
                const clearDate = new Date(clear);
                if (clearDate < enrouteDate) {
                    alert('Clear time cannot be before enroute time');
                    return;
                }
            }
            if (onScene) {
                const onSceneDate = new Date(onScene);
                const clearDate = new Date(clear);
                if (clearDate < onSceneDate) {
                    alert('Clear time cannot be before on scene time');
                    return;
                }
            }
        }
    }
    
    // Determine next step
    let nextStep = currentWizardStep + 1;
    
    // Skip step 3b (agencies) if mutual aid is not "Yes"
    if (currentWizardStep === 3 && nextStep === 3.5) {
        const mutualAidValue = document.getElementById('mutual-aid').value;
        if (mutualAidValue !== 'Yes') {
            nextStep = 4; // Skip to step 4
        }
    }
    
    // Move to next
    if (currentWizardStep < totalWizardSteps) {
        currentWizardStep = nextStep;
        updateWizardDisplay();
    }
}

function wizardPrevious() {
    if (currentWizardStep > 1) {
        let prevStep = currentWizardStep - 1;
        
        // Skip step 3b (agencies) if mutual aid is not "Yes" when going backwards
        if (currentWizardStep === 4 && prevStep === 3.5) {
            const mutualAidValue = document.getElementById('mutual-aid').value;
            if (mutualAidValue !== 'Yes') {
                prevStep = 3; // Skip back to step 3
            }
        }
        
        currentWizardStep = prevStep;
        updateWizardDisplay();
    }
}

function updateSummary(field, value) {
    const summaryMap = {
        'incident-number': 'incident',
        'call-type': 'calltype',
        'mutual-aid': 'mutualaid',
        'address': 'address',
        'town': 'town',
        'location-notes': 'locnotes',
        'dispatched': 'dispatched',
        'enroute': 'enroute',
        'on-scene': 'onscene',
        'clear': 'clear',
        'narrative': 'narrative'
    };
    
    const summaryId = summaryMap[field];
    if (summaryId) {
        const element = document.getElementById('summary-' + summaryId);
        if (element) {
            element.textContent = value || '-';
        }
    }
    
    // Highlight current field in summary
    document.querySelectorAll('.summary-item').forEach(item => {
        item.classList.remove('active');
    });
    const currentSummaryItem = document.querySelector(`.summary-item[data-summary-field="${field}"]`);
    if (currentSummaryItem) {
        currentSummaryItem.classList.add('active');
    }
}

// Add Enter key support for wizard
document.addEventListener('DOMContentLoaded', function() {
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Enter' && document.getElementById('newcall-screen').style.display !== 'none') {
            const activeElement = document.activeElement;
            if (activeElement.tagName !== 'TEXTAREA') {
                e.preventDefault();
                const btnNext = document.getElementById('btn-next');
                const btnSave = document.getElementById('btn-save');
                if (btnNext.style.display !== 'none') {
                    wizardNext();
                } else if (btnSave.style.display !== 'none') {
                    saveCall();
                }
            }
        }
    });
    
    // Add change listeners to update summary in real-time
    document.querySelectorAll('.question-card input, .question-card select, .question-card textarea').forEach(input => {
        input.addEventListener('input', function() {
            const card = this.closest('.question-card');
            const field = card.getAttribute('data-field');
            updateSummary(field, this.value);
        });
    });
});

async function saveCall() {
    // Save the narrative value from the textarea to the hidden input
    const narrativeTextarea = document.getElementById('q-narrative');
    const narrativeHidden = document.getElementById('narrative');
    if (narrativeTextarea) {
        narrativeHidden.value = narrativeTextarea.value;
    }
    
    const callType = document.getElementById('call-type').value;
    const address = document.getElementById('address').value;
    const dispatchedTime = document.getElementById('dispatched').value;
    const narrative = document.getElementById('narrative').value;
    
    if (!callType || !address || !dispatchedTime || !narrative) {
        let missing = [];
        if (!callType) missing.push('Call Type');
        if (!address) missing.push('Address');
        if (!dispatchedTime) missing.push('Dispatched Time');
        if (!narrative) missing.push('Narrative');
        alert('Please fill in all required fields. Missing: ' + missing.join(', '));
        return;
    }
    
    try {
        const dispatchedDateTime = new Date(dispatchedTime);
        
        // Collect selected apparatus
        const apparatusCheckboxes = document.querySelectorAll('input[name="apparatus"]:checked');
        const apparatusIDs = Array.from(apparatusCheckboxes).map(cb => parseInt(cb.value));
        
        // Collect selected responders
        const responderCheckboxes = document.querySelectorAll('input[name="responders"]:checked');
        const responderIDs = Array.from(responderCheckboxes).map(cb => parseInt(cb.value));
        
        const call = {
            CallType: callType,
            MutualAid: document.getElementById('mutual-aid').value,
            Address: address,
            Town: document.getElementById('town').value,
            LocationNotes: document.getElementById('location-notes').value,
            Dispatched: dispatchedDateTime.toISOString(),
            Enroute: getTimeOrNull('enroute', dispatchedDateTime),
            OnScene: getTimeOrNull('on-scene', dispatchedDateTime),
            Clear: getTimeOrNull('clear', dispatchedDateTime),
            Narrative: narrative,
            CreatedBy: currentUser.id
        };
        
        await window.go.main.App.CreateCall(call, apparatusIDs, responderIDs, []);
        alert('Call saved successfully!');
        clearNewCallForm();
        showMainMenu();
    } catch (error) {
        alert('Failed to save call: ' + error);
    }
}

function getTimeOrNull(elementId, dispatchedDateTime) {
    const value = document.getElementById(elementId).value;
    if (!value) return null;
    // Value is already in ISO format from combined date+time
    return new Date(value).toISOString();
}

function clearNewCallForm() {
    // Clear all question inputs
    document.getElementById('q-call-type').value = '';
    document.getElementById('q-mutual-aid').value = '';
    document.getElementById('q-mutual-aid-agencies-input').value = '';
    document.getElementById('q-address').value = '';
    document.getElementById('q-town').value = '';
    document.getElementById('q-location-notes').value = '';
    document.getElementById('q-dispatched-date').value = '';
    document.getElementById('q-dispatched-time').value = '';
    document.getElementById('q-enroute-date').value = '';
    document.getElementById('q-enroute-time').value = '';
    
    // Uncheck all apparatus and responders
    document.querySelectorAll('input[name="apparatus"]').forEach(cb => cb.checked = false);
    document.querySelectorAll('input[name="responders"]').forEach(cb => cb.checked = false);
    document.getElementById('q-on-scene-date').value = '';
    document.getElementById('q-on-scene-time').value = '';
    document.getElementById('q-clear-date').value = '';
    document.getElementById('q-clear-time').value = '';
    document.getElementById('q-narrative').value = '';
    
    // Clear selected agencies
    selectedAgencies = [];
    document.getElementById('selected-agencies').innerHTML = '';
    
    // Reset time to current time
    const now = new Date();
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    document.getElementById('q-dispatched-time').value = `${hours}:${minutes}`;
    
    // Clear incident display
    document.getElementById('incident-display').textContent = 'Enter date first';
    document.getElementById('incident-number-display').textContent = '-';
    document.getElementById('incident-number-header').style.display = 'none';
    
    // Clear hidden inputs
    document.getElementById('call-type').value = '';
    document.getElementById('mutual-aid').value = '';
    document.getElementById('mutual-aid-agencies').value = '';
    document.getElementById('address').value = '';
    document.getElementById('town').value = '';
    document.getElementById('location-notes').value = '';
    document.getElementById('dispatched').value = '';
    document.getElementById('enroute').value = '';
    document.getElementById('on-scene').value = '';
    document.getElementById('clear').value = '';
    document.getElementById('narrative').value = '';
    
    // Clear summary
    document.querySelectorAll('.summary-value').forEach(el => {
        el.textContent = '-';
    });
    
    // Reset to first step
    currentWizardStep = 1;
}

// Call List
async function showCallList() {
    showScreen('recent-screen');
    await loadYearSelector();
}

async function loadYearSelector() {
    const yearSelect = document.getElementById('year-selector');
    yearSelect.innerHTML = '';
    
    try {
        // Get years that have calls
        const years = await window.go.main.App.GetCallYears();
        
        if (!years || years.length === 0) {
            // No calls yet, show current year
            const currentYear = new Date().getFullYear();
            const option = document.createElement('option');
            option.value = currentYear;
            option.textContent = currentYear;
            yearSelect.appendChild(option);
        } else {
            // Add years from database
            years.forEach(year => {
                const option = document.createElement('option');
                option.value = year;
                option.textContent = year;
                yearSelect.appendChild(option);
            });
        }
        
        // Load calls for selected year
        await loadCallsByYear();
    } catch (error) {
        console.error('Failed to load years:', error);
        // Fallback to current year
        const currentYear = new Date().getFullYear();
        const option = document.createElement('option');
        option.value = currentYear;
        option.textContent = currentYear;
        yearSelect.appendChild(option);
        await loadCallsByYear();
    }
}

async function loadCallsByYear() {
    const year = parseInt(document.getElementById('year-selector').value);
    
    try {
        const calls = await window.go.main.App.GetCallsByYear(year);
        const listDiv = document.getElementById('calls-list');
        listDiv.innerHTML = '';
        
        // Handle null or undefined response
        if (!calls || calls.length === 0) {
            listDiv.innerHTML = '<p>No calls found for ' + year + '</p>';
            // Reset statistics
            document.getElementById('stat-total').textContent = '0';
            document.getElementById('stat-mutual-aid-given').textContent = '0';
            document.getElementById('stat-mutual-aid-received').textContent = '0';
            document.getElementById('stat-common-type').textContent = '-';
            return;
        }
        
        // Calculate statistics
        const stats = {
            total: calls.length,
            mutualAidGiven: calls.filter(c => c.mutual_aid === 'Yes').length,
            mutualAidReceived: calls.filter(c => c.mutual_aid === 'Received').length,
            callTypes: {}
        };
        
        calls.forEach(call => {
            stats.callTypes[call.call_type] = (stats.callTypes[call.call_type] || 0) + 1;
        });
        
        const mostCommonType = Object.keys(stats.callTypes).reduce((a, b) => 
            stats.callTypes[a] > stats.callTypes[b] ? a : b, '-');
        
        // Update statistics display
        document.getElementById('stat-total').textContent = stats.total;
        document.getElementById('stat-mutual-aid-given').textContent = stats.mutualAidGiven;
        document.getElementById('stat-mutual-aid-received').textContent = stats.mutualAidReceived;
        document.getElementById('stat-common-type').textContent = stats.total > 0 ? mostCommonType : '-';
        
        calls.forEach(call => {
            const callDiv = document.createElement('div');
            callDiv.className = 'call-item';
            callDiv.style.cursor = 'pointer';
            callDiv.onclick = () => showCallDetails(call.id);
            callDiv.innerHTML = `
                <div class="call-header">
                    <span class="call-number">${call.incident_number || 'N/A'}</span>
                    <span class="call-type">${call.call_type}</span>
                </div>
                <div class="call-details">
                    <div><strong>Address:</strong> ${call.address}, ${call.town}</div>
                    <div><strong>Dispatched:</strong> ${new Date(call.dispatched).toLocaleString()}</div>
                    <div><strong>Mutual Aid:</strong> ${call.mutual_aid}</div>
                </div>
            `;
            listDiv.appendChild(callDiv);
        });
    } catch (error) {
        console.error('Failed to load calls:', error);
        document.getElementById('calls-list').innerHTML = '<p>Error loading calls: ' + error + '</p>';
        // Reset statistics on error
        document.getElementById('stat-total').textContent = '0';
        document.getElementById('stat-mutual-aid-given').textContent = '0';
        document.getElementById('stat-mutual-aid-received').textContent = '0';
        document.getElementById('stat-common-type').textContent = '-';
    }
}

async function showCallDetails(callId) {
    try {
        const result = await window.go.main.App.GetCallByID(callId);
        const call = result.call;
        
        const modalBody = `
            <div style="text-align: left;">
                <p><strong>Incident #:</strong> ${call.incident_number || 'N/A'}</p>
                <p><strong>Call Type:</strong> ${call.call_type}</p>
                <p><strong>Address:</strong> ${call.address}, ${call.town}</p>
                <p><strong>Location Notes:</strong> ${call.location_notes || '-'}</p>
                <p><strong>Mutual Aid:</strong> ${call.mutual_aid}</p>
                <p><strong>Dispatched:</strong> ${new Date(call.dispatched).toLocaleString()}</p>
                ${call.enroute ? `<p><strong>Enroute:</strong> ${new Date(call.enroute).toLocaleString()}</p>` : ''}
                ${call.on_scene ? `<p><strong>On Scene:</strong> ${new Date(call.on_scene).toLocaleString()}</p>` : ''}
                ${call.clear ? `<p><strong>Clear:</strong> ${new Date(call.clear).toLocaleString()}</p>` : ''}
                <p><strong>Narrative:</strong></p>
                <p style="background: #f5f5f5; padding: 10px; border-radius: 4px; white-space: pre-wrap;">${call.narrative}</p>
            </div>
        `;
        
        showModal('Call Details', modalBody, 'Close');
    } catch (error) {
        alert('Failed to load call details: ' + error);
    }
}

// Placeholder functions for other screens
function showSearch() {
    showScreen('search-screen');
}

async function performSearch() {
    const query = document.getElementById('search-query').value;
    if (!query) {
        alert('Please enter a search term');
        return;
    }
    
    try {
        const calls = await window.go.main.App.SearchCalls(query);
        const resultsDiv = document.getElementById('search-results');
        resultsDiv.innerHTML = '';
        
        if (calls.length === 0) {
            resultsDiv.innerHTML = '<p>No calls found</p>';
            return;
        }
        
        calls.forEach(call => {
            const callDiv = document.createElement('div');
            callDiv.className = 'call-item';
            callDiv.innerHTML = `
                <div class="call-header">
                    <span class="call-number">${call.IncidentNumber || 'N/A'}</span>
                    <span class="call-type">${call.CallType}</span>
                </div>
                <div class="call-details">
                    <div><strong>Address:</strong> ${call.Address}</div>
                    <div><strong>Dispatched:</strong> ${new Date(call.Dispatched).toLocaleString()}</div>
                    <div><strong>Disposition:</strong> ${call.Disposition}</div>
                </div>
            `;
            resultsDiv.appendChild(callDiv);
        });
    } catch (error) {
        console.error('Search failed:', error);
        alert('Search failed: ' + error);
    }
}

function showReports() {
    showCallList(); // For now, reports = call list
}

function showExport() {
    showScreen('export-screen');
}

async function exportCSV() {
    try {
        const calls = await window.go.main.App.GetRecentCalls(1000);
        
        // Create CSV content
        let csv = 'Incident Number,Call Type,Address,Town,Dispatched,Disposition,Narrative\n';
        calls.forEach(call => {
            const row = [
                call.IncidentNumber || '',
                call.CallType,
                call.Address,
                call.Town || '',
                new Date(call.Dispatched).toLocaleString(),
                call.Disposition,
                (call.Narrative || '').replace(/"/g, '""')
            ];
            csv += row.map(field => `"${field}"`).join(',') + '\n';
        });
        
        // Download CSV
        const blob = new Blob([csv], { type: 'text/csv' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `calls-export-${new Date().toISOString().split('T')[0]}.csv`;
        a.click();
        URL.revokeObjectURL(url);
        
        document.getElementById('export-status').textContent = 'CSV exported successfully!';
        setTimeout(() => {
            document.getElementById('export-status').textContent = '';
        }, 3000);
    } catch (error) {
        alert('Export failed: ' + error);
    }
}

function exportPDF() {
    alert('PDF export coming soon!');
}

async function showAdminRoster() {
    showScreen('roster-screen');
    await loadRoster();
}

async function loadRoster() {
    try {
        const users = await window.go.main.App.GetAllUsers();
        const listDiv = document.getElementById('roster-list');
        listDiv.innerHTML = '';
        
        if (!users || users.length === 0) {
            listDiv.innerHTML = '<p>No users found</p>';
            return;
        }
        
        // Filter out the default "Admin User" (ID 0 or name "Admin User")
        const filteredUsers = users.filter(user => 
            user.id !== 0 && !(user.first_name === 'Admin' && user.last_name === 'User')
        );
        
        if (filteredUsers.length === 0) {
            listDiv.innerHTML = '<p>No users found</p>';
            return;
        }
        
        filteredUsers.forEach(user => {
            const userDiv = document.createElement('div');
            userDiv.className = 'call-item';
            // Display position with admin badge if applicable
            const positionDisplay = user.position ? (user.position.charAt(0).toUpperCase() + user.position.slice(1)) : 'Member';
            const adminBadge = user.is_admin ? ' <span style="background: #007bff; color: white; padding: 2px 8px; border-radius: 4px; font-size: 12px; margin-left: 5px;">ADMIN</span>' : '';
            const emsDisplay = user.ems_level ? ` | EMS: ${user.ems_level}` : '';
            const fullName = `${user.first_name} ${user.last_name}`;
            const joinedDateDisplay = user.joined_date ? new Date(user.joined_date).toLocaleDateString() : 'Not set';
            userDiv.innerHTML = `
                <div class="call-header">
                    <span class="call-number">${fullName}</span>
                    <span class="call-type">${positionDisplay}${adminBadge}${emsDisplay}</span>
                </div>
                <div class="call-details">
                    <div><strong>Status:</strong> ${user.active ? 'Active' : 'Inactive'}</div>
                    <div><strong>Joined:</strong> ${joinedDateDisplay}</div>
                    <div><strong>Created:</strong> ${new Date(user.created).toLocaleDateString()}</div>
                    <div style="margin-top: 10px;">
                        <button class="btn btn-secondary" onclick="editUserPosition(${user.id}, '${fullName}', '${user.position || 'member'}', '${user.ems_level || ''}', ${user.is_admin || false})">📝 Edit</button>
                        <button class="btn btn-secondary" onclick="editUserJoinDate(${user.id}, '${fullName}', '${user.joined_date || ''}')">📅 Set Join Date</button>
                        <button class="btn btn-secondary" onclick="resetUserPIN(${user.id}, '${fullName}')">🔑 Reset PIN</button>
                    </div>
                </div>
            `;
            listDiv.appendChild(userDiv);
        });
    } catch (error) {
        console.error('Failed to load roster:', error);
        document.getElementById('roster-list').innerHTML = '<p>No users found</p>';
    }
}

function showAddUser() {
    const modalBody = `
        <div class="form-group">
            <label>First Name</label>
            <input type="text" id="new-user-firstname" class="form-control">
        </div>
        <div class="form-group">
            <label>Last Name</label>
            <input type="text" id="new-user-lastname" class="form-control">
        </div>
        <div class="form-group">
            <label>Position</label>
            <input type="text" id="new-user-position" list="position-options" class="form-control" placeholder="Type or select position...">
            <datalist id="position-options">
                <option value="Chief">
                <option value="Deputy Chief">
                <option value="Captain">
                <option value="Member">
                <option value="Probationary">
            </datalist>
        </div>
        <div class="form-group">
            <label>EMS Level</label>
            <input type="text" id="new-user-ems" list="ems-options" class="form-control" placeholder="Type or select...">
            <datalist id="ems-options">
                <option value="None">
                <option value="VEFR">
                <option value="EMR">
                <option value="EMT">
                <option value="AEMT">
                <option value="Paramedic">
            </datalist>
        </div>
        <div class="form-group">
            <label>
                <input type="checkbox" id="new-user-admin"> Administrator Privileges
            </label>
        </div>
        <div class="form-group">
            <label>4-Digit PIN</label>
            <input type="password" id="new-user-pin" maxlength="4" class="form-control">
        </div>
        <div class="form-group">
            <label>Confirm PIN</label>
            <input type="password" id="new-user-pin-confirm" maxlength="4" class="form-control">
        </div>
    `;
    
    showModal('Add New User', modalBody, 'Create User', () => {
        const firstName = document.getElementById('new-user-firstname').value.trim();
        const lastName = document.getElementById('new-user-lastname').value.trim();
        const position = document.getElementById('new-user-position').value;
        const emsLevel = document.getElementById('new-user-ems').value || '';
        const isAdmin = document.getElementById('new-user-admin').checked;
        const pin = document.getElementById('new-user-pin').value;
        const pinConfirm = document.getElementById('new-user-pin-confirm').value;
        
        if (!firstName || !lastName) {
            alert('Please enter first and last name');
            return;
        }
        
        if (!position) {
            alert('Please select a position');
            return;
        }
        
        if (!pin || pin.length !== 4) {
            alert('PIN must be 4 digits');
            return;
        }
        
        if (pin !== pinConfirm) {
            alert('PINs do not match');
            return;
        }
        
        window.go.main.App.CreateUser(firstName, lastName, position, emsLevel, pin, isAdmin)
            .then(() => {
                alert('User created successfully!');
                loadRoster();
            })
            .catch(error => {
                alert('Failed to create user: ' + error);
            });
    });
}

function resetUserPIN(userID, userName) {
    const modalBody = `
        <p>Reset PIN for <strong>${userName}</strong>?</p>
        <div class="form-group">
            <label>New 4-Digit PIN</label>
            <input type="password" id="reset-pin" maxlength="4" class="form-control">
        </div>
        <div class="form-group">
            <label>Confirm PIN</label>
            <input type="password" id="reset-pin-confirm" maxlength="4" class="form-control">
        </div>
    `;
    
    showModal('Reset User PIN', modalBody, 'Reset PIN', () => {
        const newPIN = document.getElementById('reset-pin').value;
        const confirmPIN = document.getElementById('reset-pin-confirm').value;
        
        if (!newPIN || newPIN.length !== 4) {
            alert('PIN must be 4 digits');
            return;
        }
        
        if (newPIN !== confirmPIN) {
            alert('PINs do not match');
            return;
        }
        
        window.go.main.App.ChangeUserPIN(userID, newPIN)
            .then(() => {
                alert('PIN reset successfully!');
            })
            .catch(error => {
                alert('Failed to reset PIN: ' + error);
            });
    });
}

function editUserPosition(userID, userName, currentPosition, currentEMS, isAdmin) {
    const modalBody = `
        <p>Edit settings for <strong>${userName}</strong>?</p>
        <div class="form-group">
            <label>Position</label>
            <input type="text" id="edit-user-position" list="edit-position-options" class="form-control" value="${currentPosition}" placeholder="Type or select position...">
            <datalist id="edit-position-options">
                <option value="Chief">
                <option value="Deputy Chief">
                <option value="Captain">
                <option value="Member">
                <option value="Probationary">
            </datalist>
        </div>
        <div class="form-group">
            <label>EMS Level</label>
            <input type="text" id="edit-user-ems" list="edit-ems-options" class="form-control" value="${currentEMS}" placeholder="Type or select...">
            <datalist id="edit-ems-options">
                <option value="None">
                <option value="VEFR">
                <option value="EMR">
                <option value="EMT">
                <option value="AEMT">
                <option value="Paramedic">
            </datalist>
        </div>
        <div class="form-group">
            <label>
                <input type="checkbox" id="edit-user-admin" ${isAdmin ? 'checked' : ''}> Administrator Privileges
            </label>
        </div>
    `;
    
    showModal('Edit User', modalBody, 'Update', () => {
        const newPosition = document.getElementById('edit-user-position').value;
        const newEMS = document.getElementById('edit-user-ems').value;
        const newIsAdmin = document.getElementById('edit-user-admin').checked;
        
        if (!newPosition) {
            alert('Please select a position');
            return;
        }
        
        // Get the user, update fields, then save
        window.go.main.App.GetUserByID(userID)
            .then(user => {
                user.position = newPosition;
                user.ems_level = newEMS;
                user.is_admin = newIsAdmin;
                return window.go.main.App.UpdateUser(user);
            })
            .then(() => {
                alert('User updated successfully!');
                loadRoster();
            })
            .catch(error => {
                alert('Failed to update user: ' + error);
            });
    });
}

function editUserJoinDate(userID, userName, currentJoinDate) {
    const dateValue = currentJoinDate ? currentJoinDate.split('T')[0] : '';
    const modalBody = `
        <p>Set join date for <strong>${userName}</strong>?</p>
        <div class="form-group">
            <label>Join Date</label>
            <input type="date" id="edit-join-date" value="${dateValue}" class="form-control">
        </div>
    `;
    
    showModal('Set Join Date', modalBody, 'Update Date', () => {
        const joinDate = document.getElementById('edit-join-date').value;
        
        if (!joinDate) {
            alert('Please select a date');
            return;
        }
        
        window.go.main.App.UpdateUserJoinDate(userID, joinDate)
            .then(() => {
                alert('Join date updated successfully!');
                loadRoster();
            })
            .catch(error => {
                alert('Failed to update join date: ' + error);
            });
    });
}

async function showAdminPicklists() {
    showScreen('picklists-screen');
    document.getElementById('picklist-category').value = '';
    document.getElementById('picklist-items').innerHTML = '';
}

async function loadPicklistItems() {
    const category = document.getElementById('picklist-category').value;
    if (!category) return;
    
    try {
        const items = await window.go.main.App.GetPicklistByCategory(category);
        const listDiv = document.getElementById('picklist-items');
        listDiv.innerHTML = '';
        
        if (items.length === 0) {
            listDiv.innerHTML = '<p>No items found</p>';
            return;
        }
        
        items.forEach(item => {
            const itemDiv = document.createElement('div');
            itemDiv.className = 'call-item';
            itemDiv.innerHTML = `
                <div class="call-header">
                    <span class="call-number">${item.value}</span>
                    <span class="call-type">Order: ${item.sort_order}</span>
                </div>
                <div class="call-details">
                    <div><strong>Status:</strong> ${item.active ? 'Active' : 'Inactive'}</div>
                </div>
            `;
            listDiv.appendChild(itemDiv);
        });
    } catch (error) {
        console.error('Failed to load picklist items:', error);
        document.getElementById('picklist-items').innerHTML = '<p>No items found</p>';
    }
}

function showAddPicklist() {
    const category = document.getElementById('picklist-category').value;
    if (!category) {
        alert('Please select a category first');
        return;
    }
    
    const value = prompt('Enter new item value:');
    if (!value) return;
    
    const sortOrder = parseInt(prompt('Enter sort order (number):') || '99');
    
    window.go.main.App.CreatePicklist(category, value, sortOrder)
        .then(() => {
            alert('Item created successfully!');
            loadPicklistItems();
        })
        .catch(error => {
            alert('Failed to create item: ' + error);
        });
}

function showFormSettings() {
    alert('Form Settings feature coming soon!');
}

function showSettings() {
    alert('Settings feature coming soon!');
}

// Make functions globally available
window.doLogin = doLogin;
window.doLogout = doLogout;
window.showChangePIN = showChangePIN;
window.showAdminLogin = showAdminLogin;
window.showMainMenu = showMainMenu;
window.backToMenu = backToMenu;
window.showNewCall = showNewCall;
window.saveCall = saveCall;
window.showCallList = showCallList;
window.showSearch = showSearch;
window.performSearch = performSearch;
window.showReports = showReports;
window.showExport = showExport;
window.exportCSV = exportCSV;
window.exportPDF = exportPDF;
window.showAdminRoster = showAdminRoster;
window.showAddUser = showAddUser;
window.resetUserPIN = resetUserPIN;
window.showAdminPicklists = showAdminPicklists;
window.loadPicklistItems = loadPicklistItems;
window.showAddPicklist = showAddPicklist;
window.showFormSettings = showFormSettings;
window.showSettings = showSettings;
window.wizardNext = wizardNext;
window.wizardPrevious = wizardPrevious;
window.closeModal = closeModal;
window.confirmModal = confirmModal;
