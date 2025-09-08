import { useState, useEffect, useRef } from 'react';
import { Idea } from '../types';
import { getDraft, saveDraft } from '../services/storageService';

export const useAutosaveDraft = () => {
    const [ideas, setIdeas] = useState<Idea[]>([]);
    const [purpose, setPurpose] = useState('Code Generation');
    const isInitialLoad = useRef(true);

    // Load draft from storage service on mount
    useEffect(() => {
        const loadDraft = async () => {
            try {
                const savedDraft = await getDraft();
                if (savedDraft) {
                    setIdeas(savedDraft.ideas);
                    setPurpose(savedDraft.purpose);
                }
            } catch (e) {
                console.error("Failed to load draft from storage service", e);
            } finally {
                isInitialLoad.current = false;
            }
        };
        loadDraft();
    }, []);

    // Save draft to storage service on changes
    useEffect(() => {
        if (isInitialLoad.current) {
            return;
        }
        const handler = setTimeout(() => {
            try {
                saveDraft({ ideas, purpose });
            } catch (e) {
                console.error("Failed to save draft via storage service", e);
            }
        }, 500);
        return () => clearTimeout(handler);
    }, [ideas, purpose]);

    return { ideas, setIdeas, purpose, setPurpose };
};