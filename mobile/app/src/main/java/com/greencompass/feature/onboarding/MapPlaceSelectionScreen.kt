package com.greencompass.feature.onboarding

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Place
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MapPlaceSelectionScreen(
    onBack: () -> Unit,
    onConfirmPlace: () -> Unit
) {
    GreenCompassScaffold(
        title = "",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Box(modifier = Modifier.fillMaxSize().padding(paddingValues)) {
            // Simplified mock map background
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(GreenCompassColors.Mist),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    imageVector = Icons.Default.Place,
                    contentDescription = "Map placeholder",
                    tint = GreenCompassColors.ForestGreen,
                    modifier = Modifier.size(64.dp)
                )
            }

            // Bottom sheet
            Surface(
                modifier = Modifier
                    .fillMaxWidth()
                    .align(Alignment.BottomCenter),
                shape = RoundedCornerShape(topStart = 24.dp, topEnd = 24.dp),
                color = GreenCompassColors.WarmWhite,
                shadowElevation = 8.dp
            ) {
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(AppSpacing.xl)
                ) {
                    Text(
                        text = "Selected place",
                        style = GreenCompassTypography.labelMedium,
                        color = GreenCompassColors.MutedText,
                        modifier = Modifier.padding(bottom = AppSpacing.xs)
                    )

                    Text(
                        text = "Lower Valley",
                        style = GreenCompassTypography.headlineMedium,
                        color = GreenCompassColors.Charcoal,
                        modifier = Modifier.padding(bottom = AppSpacing.xl)
                    )

                    PrimaryButton(
                        text = "Confirm place",
                        onClick = onConfirmPlace
                    )
                }
            }
        }
    }
}
