package com.greencompass.feature.onboarding

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.AppSpacing
import com.greencompass.core.ui.GreenCompassColors
import com.greencompass.core.ui.GreenCompassTypography
import com.greencompass.core.ui.TextLinkButton

@Composable
fun AccountChoiceScreen(onPersonal: () -> Unit, onOrganization: () -> Unit, onInvitation: () -> Unit, onSignIn: () -> Unit) {
    val options = listOf(
        Triple("Use Green Compass personally", "Follow updates for the places that matter to you.", onPersonal),
        Triple("Join an organization", "Access information and responsibilities assigned by your organization.", onOrganization),
        Triple("I have an invitation code", "Join using an invitation sent to you.", onInvitation)
    )

    Column(modifier = Modifier.fillMaxSize().padding(horizontal = AppSpacing.lg), horizontalAlignment = Alignment.CenterHorizontally) {
        Spacer(modifier = Modifier.height(AppSpacing.xxl))
        Text(text = "How would you like to continue?", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center, modifier = Modifier.padding(bottom = AppSpacing.xl))

        LazyColumn(verticalArrangement = Arrangement.spacedBy(AppSpacing.md), modifier = Modifier.weight(1f)) {
            items(options) { (title, subtitle, onClick) ->
                Surface(
                    modifier = Modifier.fillMaxWidth().clickable { onClick() },
                    shape = RoundedCornerShape(16.dp),
                    color = Color.White,
                    border = BorderStroke(1.dp, GreenCompassColors.Stone)
                ) {
                    Column(modifier = Modifier.padding(AppSpacing.lg)) {
                        Text(text = title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xs))
                        Text(text = subtitle, style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText)
                    }
                }
            }
        }

        Spacer(modifier = Modifier.height(AppSpacing.xl))
        Text(text = "Already have an account?", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xs))
        TextLinkButton(text = "Sign in", onClick = onSignIn)
        Spacer(modifier = Modifier.height(AppSpacing.xxl))
    }
}
