// ORCA Custom JavaScript

document.addEventListener('DOMContentLoaded', function() {
  // Add external link icons
  const content = document.querySelector('.md-content');
  if (content) {
    const links = content.querySelectorAll('a[href^="http"]');
    links.forEach(link => {
      if (!link.hostname.includes('scttfrdmn.github.io')) {
        link.setAttribute('target', '_blank');
        link.setAttribute('rel', 'noopener noreferrer');
      }
    });
  }

  // Add copy button feedback
  const copyButtons = document.querySelectorAll('.md-clipboard');
  copyButtons.forEach(button => {
    button.addEventListener('click', function() {
      const originalTitle = this.getAttribute('title');
      this.setAttribute('title', 'Copied!');
      setTimeout(() => {
        this.setAttribute('title', originalTitle);
      }, 2000);
    });
  });

  // Smooth scroll for anchor links
  document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
      const target = document.querySelector(this.getAttribute('href'));
      if (target) {
        e.preventDefault();
        target.scrollIntoView({
          behavior: 'smooth',
          block: 'start'
        });
      }
    });
  });

  // Add version badge to navigation
  const version = '0.1.0-dev';
  const nav = document.querySelector('.md-header__title');
  if (nav && !document.querySelector('.version-badge')) {
    const badge = document.createElement('span');
    badge.className = 'version-badge badge badge-warning';
    badge.textContent = version;
    badge.style.marginLeft = '0.5rem';
    badge.style.fontSize = '0.7rem';
    badge.style.verticalAlign = 'middle';
    nav.appendChild(badge);
  }
});

// Analytics helper (if Google Analytics is configured)
window.dataLayer = window.dataLayer || [];
function gtag(){dataLayer.push(arguments);}
gtag('js', new Date());

// Track external link clicks
document.addEventListener('click', function(e) {
  const target = e.target.closest('a');
  if (target && target.hostname !== window.location.hostname) {
    gtag('event', 'click', {
      'event_category': 'external_link',
      'event_label': target.href,
      'transport_type': 'beacon'
    });
  }
});
