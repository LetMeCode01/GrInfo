import React, { useEffect, useState } from 'react';

export default function Leaderboard() {
  const [topUsers, setTopUsers] = useState([]);
  const [loading, setLoading] = useState(true);

  const loadLeaderboard = async () => {
    setLoading(true);

    try {
      const leaderboardResponse = await fetch('/api/leaderboard');
      if (!leaderboardResponse.ok) throw new Error('Eroare la server');
      const leaderboardData = await leaderboardResponse.json();

      const token = localStorage.getItem('token');
      if (!token) {
        setTopUsers(leaderboardData || []);
        return;
      }

      const [profileResponse, grInfoResponse] = await Promise.all([
        fetch('/api/profile', {
          headers: { Authorization: `Bearer ${token}` },
        }),
        fetch('/api/grinfo/profile', {
          headers: { Authorization: `Bearer ${token}` },
        }),
      ]);

      if (!profileResponse.ok || !grInfoResponse.ok) {
        setTopUsers(leaderboardData || []);
        return;
      }

      const profileData = await profileResponse.json();
      const grInfoProfile = await grInfoResponse.json();

      const normalizedUsers = (leaderboardData || []).map(user => {
        if (user.username !== profileData.username) return user;

        return {
          ...user,
          currentElo: Number(grInfoProfile.currentElo ?? user.currentElo ?? 1000),
        };
      });

      setTopUsers(normalizedUsers);
    } catch (err) {
      console.error('Eroare fetch:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadLeaderboard();

    const onRefresh = () => loadLeaderboard();
    window.addEventListener('focus', onRefresh);
    window.addEventListener('grinfo:profile-updated', onRefresh);

    return () => {
      window.removeEventListener('focus', onRefresh);
      window.removeEventListener('grinfo:profile-updated', onRefresh);
    };
  }, []);

  if (loading) return <div className="profile-page"><div className="profile-header"><h1>Se încarcă...</h1></div></div>;

  return (
    <div className="profile-page">
      <div className="profile-container">
        <div className="profile-header">
          <div className="profile-title">
            <h1><span className="wave-emoji">🏆</span> Leaderboard</h1>
            <p>Top utilizatori după acuratețea GrInfo</p>
          </div>
        </div>

        <div className="history-section">
          <div className="history-list">
            {topUsers.map((user, index) => (
              <div key={index} className="history-card">
                <div className="history-header">
                  <div>
                    <span className="history-course-title">
                      #{index + 1} {user.username}
                    </span>
                    <div className="history-date">
                      Corecte: {user.totalCorrectAnswers || 0} / {user.totalQuestionsAnswered || 0}
                    </div>
                  </div>
                  <div className="history-score-badge">
                    <span className="score-text good">{Number(user.accuracy || 0).toFixed(1)}%</span>
                  </div>
                </div>
                <div className="progress-track">
                  <div 
                    className="progress-fill good" 
                    style={{ width: `${Math.min(100, Number(user.accuracy || 0))}%` }}
                  ></div>
                </div>
                <div className="history-footer">ELO curent: {Number(user.currentElo || 1000).toFixed(1)}</div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}