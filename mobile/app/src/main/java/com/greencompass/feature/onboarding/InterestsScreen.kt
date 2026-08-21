package com.greencompass.feature.onboarding

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Check
import androidx.compose.material.icons.outlined.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

data class InterestOption(val title: String, val description: String, val icon: ImageVector)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun InterestsScreen(
    onBack: () -> Unit,
    onContinue: () -> Unit,
    onSkip: () -> Unit
) {
    val options = listOf(
        InterestOption("Local outlook", "Conditions for the coming days", Icons.Outlined.WbSunny),
        InterestOption("Water outlook", "Information about nearby water conditions", Icons.Outlined.WaterDrop),
        InterestOption("Seasonal information", "Changes that may affect your area", Icons.Outlined.CalendarToday),
        InterestOption("Land and ecosystems", "Changes in the surrounding environment", Icons.Outlined.Landscape),
        InterestOption("Food and agriculture", "Information relevant to the current season", Icons.Outlined.Agriculture),
        InterestOption("Community updates", "Information shared by people nearby", Icons.Outlined.People)
    )
    
    var selectedInterests by remember { mutableStateOf(setOf<String>()) }

    GreenCompassScaffold(
        title = "",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(horizontal = AppSpacing.lg)
        ) {
            Text(
                text = "Explore what matters to you",
                style = GreenCompassTypography.headlineLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.xs)
            )
            Text(
                text = "Choose a few areas to personalize your experience.\nYou can change this later.",
                style = GreenCompassTypography.bodyMedium,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            LazyColumn(
                verticalArrangement = Arrangement.spacedBy(AppSpacing.sm),
                modifier = Modifier.weight(1f)
            ) {
                items(options) { option ->
                    val isSelected = selectedInterests.contains(option.title)
                    InterestCard(
                        option = option,
                        isSelected = isSelected,
                        onClick = {
                            if (isSelected) selectedInterests -= option.title
                            else selectedInterests += option.title
                        }
                    )
                }
            }

            Spacer(modifier = Modifier.height(AppSpacing.xl))

            PrimaryButton(
                text = "Continue",
                onClick = onContinue,
                modifier = Modifier.padding(bottom = AppSpacing.sm)
            )

            TextLinkButton(
                text = "Skip for now",
                onClick = onSkip
            )

            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}

@Composable
private fun InterestCard(
    option: InterestOption,
    isSelected: Boolean,
    onClick: () -> Unit
) {
    Surface(
        modifier = Modifier
            .fillMaxWidth()
            .clickable { onClick() },
        shape = RoundedCornerShape(12.dp),
        color = if (isSelected) GreenCompassColors.SoftSage else Color.White,
        border = BorderStroke(
            width = 1.dp,
            color = if (isSelected) GreenCompassColors.ForestGreen else GreenCompassColors.Stone
        )
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(AppSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                imageVector = option.icon,
                contentDescription = null,
                tint = if (isSelected) GreenCompassColors.ForestGreen else GreenCompassColors.MutedText,
                modifier = Modifier.size(24.dp)
            )
            Spacer(modifier = Modifier.width(AppSpacing.md))
            Column(modifier = Modifier.weight(1f)) {
                Text(text = option.title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
                Text(text = option.description, style = GreenCompassTypography.bodySmall, color = GreenCompassColors.MutedText)
            }
            if (isSelected) {
                Icon(
                    imageVector = Icons.Default.Check,
                    contentDescription = "Selected",
                    tint = GreenCompassColors.ForestGreen
                )
            }
        }
    }
}
