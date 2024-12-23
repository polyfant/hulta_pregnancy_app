import React, { useState, useEffect } from 'react';
import { Box, Card, CardContent, Typography, IconButton, Collapse, useTheme, useMediaQuery } from '@mui/material';
import { styled } from '@mui/material/styles';
import {
    ExpandMore as ExpandMoreIcon,
    Female as FemaleIcon,
    Male as MaleIcon,
    Pets as PetsIcon,
    ChevronRight as ChevronRightIcon,
} from '@mui/icons-material';
import { FamilyTree as FamilyTreeType, FamilyMember } from '../../types/familyTree';

const ExpandButton = styled(IconButton)<{ expanded: boolean }>(({ expanded }) => ({
    transform: expanded ? 'rotate(90deg)' : 'rotate(0deg)',
    transition: 'transform 0.3s',
}));

const Section = styled(Box)(({ theme }) => ({
    marginBottom: theme.spacing(2),
    '&:last-child': {
        marginBottom: 0,
    },
}));

const MemberCard = styled(Card)<{ membertype: 'parent' | 'sibling' | 'offspring' }>(({ theme, membertype }) => ({
    marginBottom: theme.spacing(1),
    backgroundColor: 
        membertype === 'parent' ? theme.palette.primary.light :
        membertype === 'sibling' ? theme.palette.secondary.light :
        theme.palette.success.light,
    transition: 'transform 0.2s',
    cursor: 'pointer',
    '&:hover': {
        transform: 'translateX(4px)',
    },
}));

interface FamilyTreeProps {
    horseId: number;
    onMemberClick?: (id: number) => void;
}

export const FamilyTree: React.FC<FamilyTreeProps> = ({ horseId, onMemberClick }) => {
    const [familyData, setFamilyData] = useState<FamilyTreeType | null>(null);
    const [expandedSections, setExpandedSections] = useState({
        parents: true,
        siblings: false,
        offspring: false,
    });
    const theme = useTheme();
    const isDesktop = useMediaQuery(theme.breakpoints.up('md'));

    useEffect(() => {
        const fetchFamilyTree = async () => {
            try {
                const response = await fetch(`/api/horses/${horseId}/family`);
                if (!response.ok) throw new Error('Failed to fetch family tree');
                const data = await response.json();
                setFamilyData(data);
            } catch (error) {
                console.error('Error fetching family tree:', error);
            }
        };

        fetchFamilyTree();
    }, [horseId]);

    const toggleSection = (section: keyof typeof expandedSections) => {
        setExpandedSections(prev => ({
            ...prev,
            [section]: !prev[section],
        }));
    };

    const renderMember = (member: FamilyMember, type: 'parent' | 'sibling' | 'offspring') => {
        if (!member) return null;

        const handleClick = () => {
            if (member.id && onMemberClick) {
                onMemberClick(member.id);
            }
        };

        return (
            <MemberCard membertype={type} onClick={handleClick}>
                <CardContent>
                    <Box display="flex" alignItems="center" gap={1}>
                        {member.gender === 'MARE' ? (
                            <FemaleIcon color="primary" />
                        ) : member.gender === 'STALLION' ? (
                            <MaleIcon color="primary" />
                        ) : (
                            <PetsIcon color="primary" />
                        )}
                        <Box>
                            <Typography variant="subtitle1" component="div">
                                {member.name}
                                {member.isExternal && ' (External)'}
                            </Typography>
                            {member.breed && (
                                <Typography variant="body2" color="text.secondary">
                                    {member.breed}
                                </Typography>
                            )}
                            {member.age && (
                                <Typography variant="body2" color="text.secondary">
                                    Age: {member.age}
                                </Typography>
                            )}
                        </Box>
                    </Box>
                </CardContent>
            </MemberCard>
        );
    };

    if (!familyData) return null;

    return (
        <Box sx={{ p: 2 }}>
            <Section>
                <Box display="flex" alignItems="center" mb={1}>
                    <ExpandButton
                        expanded={expandedSections.parents}
                        onClick={() => toggleSection('parents')}
                        size="small"
                    >
                        <ChevronRightIcon />
                    </ExpandButton>
                    <Typography variant="h6">Parents</Typography>
                </Box>
                <Collapse in={expandedSections.parents}>
                    <Box sx={{ pl: 4 }}>
                        {familyData.mother && renderMember(familyData.mother, 'parent')}
                        {familyData.father && renderMember(familyData.father, 'parent')}
                        {!familyData.mother && !familyData.father && (
                            <Typography color="text.secondary">No parent information available</Typography>
                        )}
                    </Box>
                </Collapse>
            </Section>

            {familyData.siblings && familyData.siblings.length > 0 && (
                <Section>
                    <Box display="flex" alignItems="center" mb={1}>
                        <ExpandButton
                            expanded={expandedSections.siblings}
                            onClick={() => toggleSection('siblings')}
                            size="small"
                        >
                            <ChevronRightIcon />
                        </ExpandButton>
                        <Typography variant="h6">Siblings ({familyData.siblings.length})</Typography>
                    </Box>
                    <Collapse in={expandedSections.siblings}>
                        <Box sx={{ pl: 4 }}>
                            {familyData.siblings.map((sibling) => renderMember(sibling, 'sibling'))}
                        </Box>
                    </Collapse>
                </Section>
            )}

            {familyData.offspring && familyData.offspring.length > 0 && (
                <Section>
                    <Box display="flex" alignItems="center" mb={1}>
                        <ExpandButton
                            expanded={expandedSections.offspring}
                            onClick={() => toggleSection('offspring')}
                            size="small"
                        >
                            <ChevronRightIcon />
                        </ExpandButton>
                        <Typography variant="h6">Offspring ({familyData.offspring.length})</Typography>
                    </Box>
                    <Collapse in={expandedSections.offspring}>
                        <Box sx={{ pl: 4 }}>
                            {familyData.offspring.map((child) => renderMember(child, 'offspring'))}
                        </Box>
                    </Collapse>
                </Section>
            )}
        </Box>
    );
};
