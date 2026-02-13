import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { api, endpoints } from '../services/api';
import { Tweet } from '../components/tweets/Tweet';
import { FaArrowLeft, FaCalendarAlt } from 'react-icons/fa';
import { useAuth } from '../context/AuthContext';
import '../styles/profile.css';

interface UserProfile {
    id: string;
    username: string;
    firstname?: string;
    lastname?: string;
    createdAt?: string;
    description?: string; // Bio
    followersCount?: number;
    followingCount?: number;
}

interface TweetData {
    id: string;
    authorId: string;
    content: string;
    createdAt: string;
}

export const ProfilePage: React.FC = () => {
    const { userId } = useParams<{ userId: string }>();
    const navigate = useNavigate();
    const { user: currentUser } = useAuth();

    const [profile, setProfile] = useState<UserProfile | null>(null);
    const [tweets, setTweets] = useState<TweetData[]>([]);
    const [loading, setLoading] = useState(false);
    const [activeTab, setActiveTab] = useState('tweets');

    useEffect(() => {
        const fetchProfileData = async () => {
            if (!userId) return;
            setLoading(true);
            try {
                // Fetch profile
                // const profileData = await api.get(endpoints.users.profile(userId));
                // setProfile(profileData);

                // MOCK PROFILE since backend might not be fully ready with profile endpoint
                setProfile({
                    id: userId,
                    username: userId === currentUser?.id ? currentUser.username : 'user_' + userId.substring(0, 5),
                    firstname: 'John',
                    lastname: 'Doe',
                    createdAt: '2025-01-01',
                    description: 'Just a mock bio for now.',
                    followersCount: 123,
                    followingCount: 456
                });

                // Fetch tweets
                const tweetsData = await api.get(endpoints.tweets.user(userId));
                setTweets(tweetsData.tweets || []);
            } catch (err) {
                console.error('Failed to fetch profile data', err);
            } finally {
                setLoading(false);
            }
        };

        fetchProfileData();
    }, [userId, currentUser]);

    if (!profile && !loading) return <div className="p-8 text-center">User not found</div>;

    return (
        <div>
            <div className="back-header">
                <button onClick={() => navigate(-1)} className="back-button">
                    <FaArrowLeft />
                </button>
                <div>
                    <h2 className="text-xl font-bold leading-5">{profile?.username}</h2>
                    <div className="text-sm text-gray-500">{tweets.length} posts</div>
                </div>
            </div>

            <div className="profile-header">
                <div className="profile-cover"></div>
                <div className="profile-info-container">
                    <div className="profile-avatar-large">
                        {profile?.username[0]?.toUpperCase()}
                    </div>
                    <div className="profile-actions">
                        {currentUser?.id === userId ? (
                            <button className="edit-profile-btn">Edit profile</button>
                        ) : (
                            <button className="bg-white text-black font-bold px-4 py-2 rounded-full">Follow</button>
                        )}
                    </div>

                    <div className="profile-name-section">
                        <div className="profile-fullname">{profile?.firstname} {profile?.lastname}</div>
                        <div className="profile-username">@{profile?.username}</div>
                    </div>

                    <div className="profile-bio">
                        {profile?.description}
                    </div>

                    <div className="profile-metadata">
                        <div className="profile-stat">
                            <FaCalendarAlt /> Joined {profile?.createdAt ? new Date(profile.createdAt).toLocaleDateString() : 'Unknown'}
                        </div>
                    </div>

                    <div className="profile-stats-row">
                        <div>
                            <span className="stat-count">{profile?.followingCount}</span> <span className="stat-label">Following</span>
                        </div>
                        <div>
                            <span className="stat-count">{profile?.followersCount}</span> <span className="stat-label">Followers</span>
                        </div>
                    </div>
                </div>

                <div className="profile-tabs">
                    <div className={`profile-tab ${activeTab === 'tweets' ? 'active' : ''}`} onClick={() => setActiveTab('tweets')}>
                        Posts
                        {activeTab === 'tweets' && <div className="tab-indicator" />}
                    </div>
                    <div className={`profile-tab ${activeTab === 'likes' ? 'active' : ''}`} onClick={() => setActiveTab('likes')}>
                        Likes
                        {activeTab === 'likes' && <div className="tab-indicator" />}
                    </div>
                </div>
            </div>

            <div className="profile-content">
                {loading ? (
                    <div className="p-8 text-center text-gray-500">Loading...</div>
                ) : (
                    tweets.length > 0 ? (
                        tweets.map(tweet => <Tweet key={tweet.id} tweet={tweet} />)
                    ) : (
                        <div className="p-8 text-center text-gray-500">No tweets to show</div>
                    )
                )}
            </div>
        </div>
    );
};
