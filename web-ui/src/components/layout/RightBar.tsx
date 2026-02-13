import React from 'react';
import { FaSearch } from 'react-icons/fa';
import '../../styles/layout.css';

export const RightBar: React.FC = () => {
    return (
        <aside className="layout-rightbar">
            <div className="search-container">
                <FaSearch className="text-gray-500" />
                <input
                    type="text"
                    placeholder="Search Twitter"
                    className="search-input"
                />
            </div>

            <div className="widget">
                <h3 className="widget-header">What's happening</h3>

                {[1, 2, 3, 4].map((i) => (
                    <div key={i} className="widget-item">
                        <div className="trend-meta">
                            <span>Technology · Trending</span>
                            <span>...</span>
                        </div>
                        <div className="trend-name">#ReactJs</div>
                        <div className="trend-meta">125K Tweets</div>
                    </div>
                ))}

                <div className="show-more">Show more</div>
            </div>

            <div className="widget">
                <h3 className="widget-header">Who to follow</h3>

                {[1, 2].map((i) => (
                    <div key={i} className="widget-item flex justify-between items-center">
                        <div className="flex gap-3">
                            <div className="avatar" style={{ width: 32, height: 32 }}></div>
                            <div>
                                <div className="font-bold">Jane Doe</div>
                                <div className="trend-meta">@janedoe</div>
                            </div>
                        </div>
                        <button className="follow-btn">
                            Follow
                        </button>
                    </div>
                ))}

                <div className="show-more">Show more</div>
            </div>
        </aside>
    );
};
