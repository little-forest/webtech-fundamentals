document.addEventListener("DOMContentLoaded", function () {
  const calendarContainer = document.getElementById("calendar");
  const today = new Date();
  let currentMonth = today.getMonth();
  let currentYear = today.getFullYear();

  function generateCalendar(month, year) {
    const calendarDate = new Date(year, month, 1);
    const firstDay = calendarDate.getDay();
    const daysInMonth = new Date(year, month + 1, 0).getDate();
    const monthNames = [
      "January", "February", "March", "April",
      "May", "June", "July", "August",
      "September", "October", "November", "December"
    ];

    const table = document.createElement("table");
    const caption = table.createCaption();
    caption.textContent = `${monthNames[month]} ${year}`;

    const thead = table.createTHead();
    const tr = thead.insertRow();

    // Create table headers for days of the week
    const daysOfWeek = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
    for (const day of daysOfWeek) {
      const th = document.createElement("th");
      th.textContent = day;
      tr.appendChild(th);
    }

    const tbody = table.createTBody();
    let date = 1;

    // Create calendar cells
    for (let i = 0; i < 6; i++) {
      const row = tbody.insertRow();

      for (let j = 0; j < 7; j++) {
        if (i === 0 && j < firstDay) {
          const cell = row.insertCell();
          cell.textContent = "";
        } else if (date > daysInMonth) {
          break;
        } else {
          const cell = row.insertCell();
          cell.textContent = date;
          if (date === today.getDate() && year === today.getFullYear() && month === today.getMonth()) {
            cell.classList.add("today");
          }
          date++;
        }
      }
    }

    calendarContainer.innerHTML = "";
    calendarContainer.appendChild(table);
  }

  generateCalendar(currentMonth, currentYear);
});

