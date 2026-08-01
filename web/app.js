// Theme Toggle Controller
const themeToggleBtn = document.getElementById("theme-toggle");
const rootHtml = document.documentElement;

// Initialize theme from storage or system preference
const savedTheme = localStorage.getItem("theme");
const systemPrefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
const initialTheme = savedTheme || (systemPrefersDark ? "dark" : "light");

setTheme(initialTheme);

themeToggleBtn.addEventListener("click", () => {
    const currentTheme = rootHtml.getAttribute("data-theme");
    const newTheme = currentTheme === "dark" ? "light" : "dark";
    setTheme(newTheme);
});

function setTheme(theme) {
    rootHtml.setAttribute("data-theme", theme);
    localStorage.setItem("theme", theme);
    
    // Update Toggle Icon
    const icon = themeToggleBtn.querySelector("i");
    if (theme === "dark") {
        icon.className = "fa-solid fa-sun";
    } else {
        icon.className = "fa-solid fa-moon";
    }
}

// Quick fill chip handler
function fillAndSearch(domain) {
    const input = document.getElementById("domain-input");
    input.value = domain;
    performLookup(domain);
}

// Handle Form Submission
function handleFormSubmit(event) {
    event.preventDefault();
    const domainInput = document.getElementById("domain-input");
    const domain = domainInput.value.trim().toLowerCase();
    
    if (domain) {
        performLookup(domain);
    }
}

// Perform Lookup with Raccoon Scanning Steps
async function performLookup(domain) {
    const loader = document.getElementById("loader");
    const errorCard = document.getElementById("error-card");
    const resultsCard = document.getElementById("results-card");
    const verifyBtn = document.getElementById("verify-btn");
    
    const stepMx = document.getElementById("step-mx");
    const stepSpf = document.getElementById("step-spf");
    const stepDmarc = document.getElementById("step-dmarc");

    // Reset views
    errorCard.classList.add("hidden");
    resultsCard.classList.add("hidden");
    
    // Reset steps state
    resetSteps([stepMx, stepSpf, stepDmarc]);
    
    // Start Loading State
    loader.classList.remove("hidden");
    verifyBtn.disabled = true;

    // Animated Step Progress Simulation
    stepMx.classList.add("active");
    
    const apiPromise = fetch(`/api/verify?domain=${encodeURIComponent(domain)}`);
    
    setTimeout(() => {
        markStepDone(stepMx);
        stepSpf.classList.add("active");
    }, 400);

    setTimeout(() => {
        markStepDone(stepSpf);
        stepDmarc.classList.add("active");
    }, 800);

    try {
        const response = await apiPromise;
        if (!response.ok) {
            const errData = await response.json();
            throw new Error(errData.error || `HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        
        // Complete remaining steps
        markStepDone(stepDmarc);
        
        // Brief delay for visual completion
        setTimeout(() => {
            loader.classList.add("hidden");
            displayResults(data);
            verifyBtn.disabled = false;
        }, 400);

    } catch (error) {
        console.error("Verification error:", error);
        loader.classList.add("hidden");
        document.getElementById("error-message").textContent = error.message || "The Raccoon could not complete the domain search.";
        errorCard.classList.remove("hidden");
        verifyBtn.disabled = false;
    }
}

function resetSteps(steps) {
    steps.forEach(step => {
        step.className = "step-item";
        const icon = step.querySelector(".step-icon");
        icon.className = "fa-solid fa-circle-notch fa-spin step-icon";
    });
}

function markStepDone(step) {
    step.className = "step-item done";
    const icon = step.querySelector(".step-icon");
    icon.className = "fa-solid fa-circle-check step-icon";
}

function displayResults(data) {
    const resultsCard = document.getElementById("results-card");
    const resultDomain = document.getElementById("result-domain");
    const overallBadge = document.getElementById("overall-status-badge");
    const overallText = document.getElementById("overall-status-text");
    const overallIcon = document.getElementById("overall-status-icon");
    
    // Domain header
    resultDomain.textContent = data.domain;
    
    // Overall status
    if (data.is_valid) {
        overallText.textContent = "Raccoon Verified";
        overallIcon.className = "fa-solid fa-circle-check";
        overallBadge.className = "status-badge status-passed";
    } else {
        overallText.textContent = "Raccoon Flagged";
        overallIcon.className = "fa-solid fa-circle-exclamation";
        overallBadge.className = "status-badge status-failed";
    }
    
    // Record Rows
    updateRecordRow(
        "mx", 
        data.has_mx, 
        data.has_mx ? "Active mail exchange servers found" : "No MX records found (domain cannot receive emails)"
    );
    
    updateRecordRow(
        "spf", 
        data.has_spf, 
        data.has_spf ? "Sender Policy Framework is active" : "No SPF record found (vulnerable to domain spoofing)"
    );
    
    updateRecordRow(
        "dmarc", 
        data.has_dmarc, 
        data.has_dmarc ? "DMARC policy enforcement is active" : "No DMARC record found (vulnerable to phishing attacks)"
    );
    
    // SPF Content & Tip
    const spfBlock = document.getElementById("spf-record-block");
    const spfRaw = document.getElementById("spf-raw");
    const spfTip = document.getElementById("spf-tip-block");
    
    if (data.has_spf && data.spf_record) {
        spfRaw.textContent = data.spf_record;
        spfBlock.classList.remove("hidden");
        spfTip.classList.add("hidden");
    } else {
        spfBlock.classList.add("hidden");
        spfTip.classList.remove("hidden");
    }
    
    // DMARC Content & Tip
    const dmarcBlock = document.getElementById("dmarc-record-block");
    const dmarcRaw = document.getElementById("dmarc-raw");
    const dmarcTip = document.getElementById("dmarc-tip-block");
    
    if (data.has_dmarc && data.dmarc_record) {
        dmarcRaw.textContent = data.dmarc_record;
        dmarcBlock.classList.remove("hidden");
        dmarcTip.classList.add("hidden");
    } else {
        dmarcBlock.classList.add("hidden");
        dmarcTip.classList.remove("hidden");
    }
    
    // Reveal Results Card with smooth scroll
    resultsCard.classList.remove("hidden");
    resultsCard.scrollIntoView({ behavior: "smooth", block: "nearest" });
}

function updateRecordRow(prefix, exists, message) {
    const card = document.getElementById(`${prefix}-card`);
    const icon = document.getElementById(`${prefix}-icon`);
    const desc = document.getElementById(`${prefix}-desc`);
    
    desc.textContent = message;
    
    if (exists) {
        card.className = "status-row-card status-success";
        icon.className = "fa-solid fa-circle-check";
    } else {
        card.className = "status-row-card status-error";
        icon.className = "fa-solid fa-circle-xmark";
    }
}
