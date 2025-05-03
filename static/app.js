const API_URL = "http://localhost:8080";

let habits = [];
let nextId = 1;

const habitForm = document.getElementById("habit-form");
const habitInput = document.getElementById("habit-name");
const habitList = document.getElementById("habit-list");

window.addEventListener("DOMContentLoaded", loadHabits);

function loadHabits() {
    fetch(API_URL)
        .then((res) => {
            console.log("resp status: ",res.status);
            return res.json()
        })
        .then((data) => {
            habits = data;
            console.log("data from backend:", data);
            renderHabits();
            loadHeatmap();
        })
        .catch((err) => console.error("Error during load:", err));
}  

habitForm.addEventListener("submit", function (e) {
    e.preventDefault();
    const name = habitInput.value.trim();
    if (!name) return;
    const newhabit = {
        name,
        progress: 0,
        completedToday: false
    };
    fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify( newhabit ),
    })
        .then((res) => res.json())
        .then((newHabit) => {
            habits.push(newHabit);
            renderHabits();
            habitInput.value = "";
            loadHeatmap();
        })
        .catch((err) => console.error("Error during adding:", err));
});

function deleteHabit(id) {
    fetch(`${API_URL}/${id}`, {
        method: "DELETE",
    })
        .then(() => {
            habits = habits.filter((h) => h.id !== id);
            renderHabits();
            loadHeatmap();
        })
        .catch((err) => console.error("Ошибка при удалении:", err));
}

function completeHabit(id) {
    const habit = habits.find((h) => h.id === id);
    if (!habit) return;
    let delta = habit.completedToday !== true ? 1 : -1;
    const today = new Date().toISOString().split("T")[0];
    fetch(`${API_URL}/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: habit.name , progress: habit.progress + delta, completedToday: (habit.completedToday!==true) , date: today}),
    })
    .then((res) => res.json())
    .then((updatedHabit) => {
        const index = habits.findIndex((h) => h.id === id);
        habits[index] = updatedHabit;
        renderHabits();
        loadHeatmap();
    })
    .catch((err) => console.error("Error updating habit:", err));
}  

function renderHabits() {
    habitList.innerHTML = "";

    habits.forEach((habit) => {
        const li = document.createElement("li");

        const nameSpan = document.createElement("span");
        nameSpan.textContent = habit.name;
        li.appendChild(nameSpan);

        const progressSpan = document.createElement("span");
        progressSpan.textContent = ` (${habit.progress}/21 days) `;
        progressSpan.style.marginLeft = "10px";
        li.appendChild(progressSpan);

        const checkbox = document.createElement("input");
        checkbox.type = "checkbox";
        checkbox.checked = habit.completedToday;
        checkbox.className = "habit-checkbox";
        checkbox.onchange = () => completeHabit(habit.id);
        li.appendChild(checkbox);

        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete";
        deleteBtn.style.marginLeft = "10px";
        deleteBtn.onclick = () => deleteHabit(habit.id);
        li.appendChild(deleteBtn);

        habitList.appendChild(li);
    });
}

function generateHeatmap(dataByDate) {
    const heatmap = document.getElementById("heatmap");
    heatmap.innerHTML = "";

    const start = new Date();
    start.setFullYear(start.getFullYear() - 1);
    const totalDays = 365;
    const weeks = Math.ceil(totalDays / 7);

    for (let w = 0; w < weeks; w++) {
        const week = document.createElement("div");
        week.className = "week";

        for (let d = 0; d < 7; d++) {
            const day = document.createElement("div");
            day.className = "day";

            const date = new Date(start);
            date.setDate(start.getDate() + w * 7 + d);
            const iso = date.toISOString().split("T")[0];
            const count = dataByDate[iso] || 0;
            if (count > 0) {
                day.classList.add(`level-${Math.min(3, count)}`);
            }

            week.appendChild(day);
        }

        heatmap.appendChild(week);
    }
}

function loadHeatmap() {
    fetch(`${API_URL}/api/heatmap`)
      .then(res => res.json())
      .then(data => generateHeatmap(data))
      .catch(err => console.error("Heatmap load failed:", err));
}
  