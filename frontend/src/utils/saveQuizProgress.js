// Helper function to save quiz progress to backend
export async function saveQuizProgress(courseId, correctAnswers, totalQuestions) {
  const token = localStorage.getItem("token");
  if (!token) {
    console.warn("⚠️ Not authenticated, quiz progress not saved");
    return null;
  }

  try {
    console.log(`📤 Saving quiz progress: ${courseId} - ${correctAnswers}/${totalQuestions}`);
    
    const response = await fetch("/api/progress", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`,
      },
      body: JSON.stringify({
        courseId,
        quizScore: correctAnswers,
        totalQuestions,
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const data = await response.json();
    console.log(`✅ Quiz progress saved! Earned ${data.xpEarned} XP`);
    return data;
  } catch (error) {
    console.error("❌ Error saving quiz progress:", error);
    return null;
  }
}
